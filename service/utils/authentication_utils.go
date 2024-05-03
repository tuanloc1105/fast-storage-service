package utils

import (
	"context"
	"errors"
	"fast-storage-go-service/constant"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type TokenInformation struct {
	Exp               int64          `json:"exp"`
	Iat               int64          `json:"iat"`
	Jti               string         `json:"jti"`
	Iss               string         `json:"iss"`
	Aud               []string       `json:"aud"`
	Sub               string         `json:"sub"`
	Typ               string         `json:"typ"`
	Azp               string         `json:"azp"`
	SessionState      string         `json:"session_state"`
	ACR               string         `json:"acr"`
	AllowedOrigins    []string       `json:"allowed-origins"`
	RealmAccess       RealmAccess    `json:"realm_access"`
	ResourceAccess    ResourceAccess `json:"resource_access"`
	Scope             string         `json:"scope"`
	Sid               string         `json:"sid"`
	EmailVerified     bool           `json:"email_verified"`
	PreferredUsername string         `json:"preferred_username"`
}

type RealmAccess struct {
	Roles []string `json:"roles"`
}

type ResourceAccess struct {
	MasterRealm RealmAccess `json:"master-realm"`
	Account     RealmAccess `json:"account"`
}

func EncryptPassword(password string) (encryptedPassword string, error error) {
	encryptedPassword = ""
	bytePassword := []byte(password)
	hashedPassword, generateFromPasswordErr := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if generateFromPasswordErr == nil {
		encryptedPassword = string(hashedPassword)
	} else {
		error = generateFromPasswordErr
	}
	return encryptedPassword, error
}

func EncryptPasswordPointer(password *string) (error error) {
	bytePassword := []byte(*password)
	hashedPassword, generateFromPasswordErr := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if generateFromPasswordErr == nil {
		*password = string(hashedPassword)
	}
	error = generateFromPasswordErr
	return error
}

func VerifyJwtToken(ctx context.Context, token string) (TokenInformation, error) {
	result := TokenInformation{}

	tokenArray := strings.Split(token, ".")

	if len(tokenArray) < 3 {
		return result, errors.New("invalid token")
	}

	tokenBody := tokenArray[1] // + "="

	tokenOut, _, tokenError := Shellout(ctx, fmt.Sprintf("echo '%s' | base64 -d", tokenBody))

	if tokenError != nil && tokenOut == "" {
		return result, tokenError
	}

	usernameFromContext := ctx.Value(constant.UsernameLogKey)
	traceIdFromContext := ctx.Value(constant.TraceIdLogKey)
	username := ""
	traceId := ""
	if usernameFromContext != nil {
		username = usernameFromContext.(string)
	}
	if traceIdFromContext != nil {
		traceId = traceIdFromContext.(string)
	}

	JsonToStruct(tokenOut, &result)

	if result.Jti == "" {
		log.Error(
			fmt.Sprintf(
				constant.LogPattern,
				traceId,
				username,
				"can not parse token",
			),
		)
		return result, errors.New("can not parse token")
	}

	// validate if token is expired

	currentTimeUnix := time.Now().Unix()

	if result.Exp < currentTimeUnix {
		return result, errors.New("token expired")
	}

	return result, nil
}

func ComparePassword(inputPassword string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
}

func GetCurrentUsername(c *gin.Context) (username *string, err error) {

	currentUser, isCurrentUserExist := c.Get("auth")

	emptyString := constant.EmptyString

	if !isCurrentUserExist {
		return &emptyString, errors.New("can not get current username")
	}

	claim := currentUser.(TokenInformation)

	currentUsername := claim.PreferredUsername

	return &currentUsername, nil
}
