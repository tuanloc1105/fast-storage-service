package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fast-storage-go-service/constant"
	"fast-storage-go-service/log"
	"fast-storage-go-service/payload"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Permission struct {
	Url  string   `json:"url"`
	Role []string `json:"role"`
}

type BodyLogWriter struct {
	// embedded field. This concept allows one struct to include another struct (or an interface)
	// as if it were a field, but with additional benefits such as method promotion and type embedding.
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w BodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func ErrorHandler(c *gin.Context) {
	CheckAndSetTraceId(c)
	if c.Errors != nil && len(c.Errors.Errors()) != 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": c.Errors.Errors()})
	}
}

func RequestLogger(c *gin.Context) {
	CheckAndSetTraceId(c)
	var buf bytes.Buffer
	tee := io.TeeReader(c.Request.Body, &buf)
	body, _ := io.ReadAll(tee)
	c.Request.Body = io.NopCloser(&buf)
	dst := &bytes.Buffer{}
	isPossibleToLogRequestPayload := true
	if err := json.Compact(dst, body); err != nil && len(body) > 0 {
		isPossibleToLogRequestPayload = false
	}

	header := map[string][]string(c.Request.Header)

	headerString := ""

	for k, v := range header {
		if IsSensitiveField(k) {
			headerString += fmt.Sprintf("\n\t\t+ %s: %s", k, "***")
		} else {
			headerString += fmt.Sprintf("\n\t\t+ %s: %s", k, strings.Join(v, ", "))
		}
	}
	token := c.Request.Header.Get("Authorization")
	var requestUri string
	if token != "" {
		requestUri = strings.Replace(c.Request.RequestURI, token, "***", -1)
	} else {
		requestUri = c.Request.RequestURI
	}

	xForwardedFor := c.Request.Header.Get("X-Forwarded-For")
	xRealIP := c.Request.Header.Get("X-Real-IP")

	var message string
	if isPossibleToLogRequestPayload {
		message = fmt.Sprintf(
			"Request info:\n\t- header:%s\n\t- url: %s\n\t- X-Forwarded-For: %s\n\t- X-Real-IP: %s\n\t- ClientIP: %s\n\t- method: %s\n\t- proto: %s\n\t- payload:\n\t%s",
			headerString,
			requestUri,
			xForwardedFor,
			xRealIP,
			c.ClientIP(),
			c.Request.Method,
			c.Request.Proto,
			dst.String(),
		)
	} else {
		message = fmt.Sprintf(
			"Request info:\n\t- header:%s\n\t- url: %s\n\t- X-Forwarded-For: %s\n\t- X-Real-IP: %s\n\t- ClientIP: %s\n\t- method: %s\n\t- proto: %s",
			headerString,
			requestUri,
			xForwardedFor,
			xRealIP,
			c.ClientIP(),
			c.Request.Method,
			c.Request.Proto,
		)
	}

	currentUser := "unknown"
	claimFromGinContext, _ := c.Get("auth")
	if claimFromGinContext != nil {
		claims := claimFromGinContext.(TokenInformation)
		currentUser = claims.PreferredUsername
	}
	var ctx = context.Background()
	ctx = context.WithValue(ctx, constant.UsernameLogKey, currentUser)
	ctx = context.WithValue(ctx, constant.TraceIdLogKey, GetTraceId(c))
	log.WithLevel(
		constant.Info,
		ctx,
		HideSensitiveJsonField(message),
	)
	c.Next()
	// latency := time.Since(t)
	// log.Info("%s %s %s %s\n",
	// 	c.Request.RequestURI,
	// )
}

func ResponseLogger(c *gin.Context) {
	startingTime := time.Now()
	CheckAndSetTraceId(c)
	// c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
	// c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	// c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, PUT, GET, OPTIONS, DELETE")
	// c.Writer.Header().Set("Access-Control-Max-Age", "3600")
	// c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Content-Length, X-Requested-With, credential, X-XSRF-TOKEN")
	blw := &BodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw

	c.Next()

	latency := time.Since(startingTime)

	header := map[string][]string(c.Writer.Header())

	headerString := ""

	for k, v := range header {
		if IsSensitiveField(k) {
			headerString += fmt.Sprintf("\n\t\t+ %s: %s", k, "***")
		} else {
			headerString += fmt.Sprintf("\n\t\t+ %s: %s", k, strings.Join(v, ", "))
		}
	}
	// (*blw).body = bytes.NewBufferString("") <=> blw.body = bytes.NewBufferString("")
	responseSizeMB := bytesToMB(int32((*blw).Size()))
	responseSizeKB := bytesToKB(int32((*blw).Size()))
	statusCode := c.Writer.Status()
	var message string
	token := c.Request.Header.Get("Authorization")
	var requestUri string
	if token != "" {
		requestUri = strings.Replace(c.Request.RequestURI, token, "***", -1)
	} else {
		requestUri = c.Request.RequestURI
	}
	if responseSizeKB > float64(10) {
		message = fmt.Sprintf(
			"Response info:\n\t- status code: %s\n\t- method: %s\n\t- url: %s\n\t- header:%s\n\t- response size: %.6f MB\n\t- latency: %s",
			strconv.Itoa(statusCode),
			c.Request.Method,
			requestUri,
			headerString,
			responseSizeMB,
			latency,
		)
	} else {
		message = fmt.Sprintf(
			"Response info:\n\t- status code: %s\n\t- method: %s\n\t- url: %s\n\t- header:%s\n\t- response size: %.6f MB\n\t- latency: %s\n\t- payload:\n\t%s",
			strconv.Itoa(statusCode),
			c.Request.Method,
			requestUri,
			headerString,
			responseSizeMB,
			latency,
			blw.body.String(),
		)
	}
	currentUser := "unknown"
	claimFromGinContext, _ := c.Get("auth")
	if claimFromGinContext != nil {
		claims := claimFromGinContext.(TokenInformation)
		currentUser = claims.PreferredUsername
	}
	var ctx = context.Background()
	ctx = context.WithValue(ctx, constant.UsernameLogKey, currentUser)
	ctx = context.WithValue(ctx, constant.TraceIdLogKey, GetTraceId(c))
	log.WithLevel(
		constant.Info,
		ctx,
		HideSensitiveJsonField(message),
	)

}

func GetTokenInParamAndSetToHeader() func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenToParam := c.Query("token")
		c.Request.Header.Set("Authorization", tokenToParam)
		c.Next()
	}
}

func AuthenticationWithAuthorization(listOfRole []string) func(c *gin.Context) {
	return func(c *gin.Context) {
		CheckAndSetTraceId(c)
		traceId := GetTraceId(c)
		ctx := context.Background()
		ctx = context.WithValue(ctx, constant.TraceIdLogKey, traceId)
		token := c.Request.Header.Get("Authorization")
		var mapClaims TokenInformation
		var verifyJwtTokenError error
		if strings.Contains(token, "Bearer") {
			mapClaims, verifyJwtTokenError = VerifyJwtToken(ctx, token[7:])
		} else {
			mapClaims, verifyJwtTokenError = VerifyJwtToken(ctx, token)
		}
		if verifyJwtTokenError != nil {
			log.WithLevel(
				constant.Error,
				ctx,
				"token invalid: %s",
				verifyJwtTokenError.Error(),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, &payload.Response{
				Trace:        traceId,
				ErrorCode:    constant.Unauthorized.ErrorCode,
				ErrorMessage: constant.Unauthorized.ErrorMessage + ". " + verifyJwtTokenError.Error(),
			})
			return
		}
		currentUsername := mapClaims.PreferredUsername
		ctx = context.WithValue(ctx, constant.UsernameLogKey, currentUsername)
		c.Set("auth", mapClaims)
		var requestUri string
		if token != "" {
			requestUri = strings.Replace(c.Request.RequestURI, token, "***", -1)
		} else {
			requestUri = c.Request.RequestURI
		}
		log.WithLevel(
			constant.Info,
			ctx,
			fmt.Sprintf("Check permission for url: %v", requestUri),
		)
		if listOfRole == nil || len(listOfRole) < 1 {
			c.Next()
			return
		}
		userRolesFromAccessToken := mapClaims.RealmAccess.Roles
		if userRolesFromAccessToken != nil {
			log.WithLevel(
				constant.Info,
				ctx,
				fmt.Sprintf(
					"\n\t- this user has role: %v\n\t- current api require user with role: %v",
					userRolesFromAccessToken,
					listOfRole,
				),
			)
			for _, roleElement := range userRolesFromAccessToken {
				if slices.Contains(listOfRole, fmt.Sprintf("%v", roleElement)) {
					c.Next()
					return
				}
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, &payload.Response{
			Trace:        traceId,
			ErrorCode:    constant.Forbidden.ErrorCode,
			ErrorMessage: constant.Forbidden.ErrorMessage,
		})
	}
}

func ReturnResponse(c *gin.Context, errCode constant.ErrorEnums, responseData any, additionalMessage ...string) *payload.Response {
	message := ""
	if len(additionalMessage) > 0 {
		message = additionalMessage[0]
	}

	return &payload.Response{
		Trace:        GetTraceId(c),
		ErrorCode:    errCode.ErrorCode,
		ErrorMessage: strings.Replace(strings.Trim(strings.Trim(errCode.ErrorMessage, " ")+". "+strings.Trim(message, " ")+".", " "), ". .", ".", -1),
		Response:     responseData,
	}
}

// func ReturnPageResponse(
// 	c *gin.Context,
// 	errCode constant.ErrorEnums,
// 	totalElement int64,
// 	totalPage int64,
// 	responseData any,
// 	additionalMessage ...string,
// ) *payload.PageResponse {
// 	message := ""
// 	if len(additionalMessage) > 0 {
// 		message = additionalMessage[0]
// 	}

// 	return &payload.PageResponse{
// 		Trace:        GetTraceId(c),
// 		ErrorCode:    errCode.ErrorCode,
// 		ErrorMessage: strings.Replace(strings.Trim(strings.Trim(errCode.ErrorMessage, " ")+". "+strings.Trim(message, " ")+".", " "), ". .", ".", -1),
// 		TotalElement: totalElement,
// 		TotalPage:    totalPage,
// 		Response:     responseData,
// 	}
// }

func ReadGinContextToPayload[T any](c *gin.Context, requestPayload *T) bool {
	if err := c.ShouldBindJSON(requestPayload); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			ReturnResponse(
				c,
				constant.JsonBindingError,
				nil,
				err.Error(),
			),
		)
		return false
	}
	return true
}

func bytesToMB(bytes int32) float64 {
	result := float64(bytes) / float64(constant.BytesInMB)
	return result
}

func bytesToKB(bytes int32) float64 {
	result := float64(bytes) / float64(constant.BytesInKB)
	return result
}
