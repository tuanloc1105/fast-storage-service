package services

import (
	"context"
	"errors"
	"fast-storage-go-service/constant"
	"fast-storage-go-service/keycloak"
	"fast-storage-go-service/model"
	"fast-storage-go-service/payload"
	"fast-storage-go-service/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	KeycloakGrantType = "password"
	KeycloakScope     = "openid"
)

type AuthenticateHandler struct {
	DB  *gorm.DB
	Ctx context.Context
}

func (h AuthenticateHandler) Login(c *gin.Context) {

	ctx, isSuccess := utils.PrepareContext(c, true)
	if !isSuccess {
		return
	}
	h.Ctx = ctx

	requestPayload := payload.LoginRequestBody{}
	isParseRequestPayloadSuccess := utils.ReadGinContextToPayload(c, &requestPayload)
	if !isParseRequestPayloadSuccess {
		return
	}

	if loginResult, loginError := keycloak.KeycloakLogin(h.Ctx, requestPayload.Request.Username, requestPayload.Request.Password); loginError != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.AuthenticateFailure,
				nil,
				loginError.Error(),
			),
		)
	} else {
		if loginResult.Error != "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				utils.ReturnResponse(
					c,
					constant.AuthenticateFailure,
					loginResult,
				),
			)
			return
		}
		protocolOpenidConnectTokenResponse := payload.ProtocolOpenidConnectTokenResponse(loginResult)
		c.JSON(
			http.StatusOK,
			utils.ReturnResponse(
				c,
				constant.Success,
				protocolOpenidConnectTokenResponse,
			),
		)
	}
}

func (h AuthenticateHandler) GetUserInfo(c *gin.Context) {

	ctx, isSuccess := utils.PrepareContext(c)
	if !isSuccess {
		return
	}
	h.Ctx = ctx

	if userInfoResult, userInfoError := keycloak.KeycloakGetUserInfo(h.Ctx, c.GetHeader("Authorization")[7:]); userInfoError != nil {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			utils.ReturnResponse(
				c,
				constant.AuthenticateFailure,
				nil,
				userInfoError.Error(),
			),
		)
	} else {
		if !userInfoResult.Active {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				utils.ReturnResponse(
					c,
					constant.Unauthorized,
					nil,
				),
			)
			return
		}
		realmAccessResponse := payload.RealmAccessResponse{
			Roles: userInfoResult.RealmAccess.Roles,
		}
		resourceAccessResponse := payload.ResourceAccessResponse{
			MasterRealm: payload.RealmAccessResponse(userInfoResult.ResourceAccess.MasterRealm),
			Account:     payload.RealmAccessResponse(userInfoResult.ResourceAccess.Account),
		}
		openidConnectTokenIntrospectResponse := payload.OpenidConnectTokenIntrospectResponse{
			Exp:               userInfoResult.Exp,
			Iat:               userInfoResult.Iat,
			Jti:               userInfoResult.Jti,
			Iss:               userInfoResult.Iss,
			Aud:               userInfoResult.Aud,
			Sub:               userInfoResult.Sub,
			Typ:               userInfoResult.Typ,
			Azp:               userInfoResult.Azp,
			SessionState:      userInfoResult.SessionState,
			ACR:               userInfoResult.ACR,
			AllowedOrigins:    userInfoResult.AllowedOrigins,
			RealmAccess:       realmAccessResponse,
			ResourceAccess:    resourceAccessResponse,
			Scope:             userInfoResult.Scope,
			Sid:               userInfoResult.Sid,
			EmailVerified:     userInfoResult.EmailVerified,
			PreferredUsername: userInfoResult.PreferredUsername,
			ClientID:          userInfoResult.ClientID,
			Username:          userInfoResult.Username,
			TokenType:         userInfoResult.TokenType,
			Active:            userInfoResult.Active,
		}
		c.JSON(
			http.StatusOK,
			utils.ReturnResponse(
				c,
				constant.Success,
				openidConnectTokenIntrospectResponse,
			),
		)
	}
}

func (h AuthenticateHandler) GetNewToken(c *gin.Context) {

	ctx, isSuccess := utils.PrepareContext(c)
	if !isSuccess {
		return
	}
	h.Ctx = ctx

	requestPayload := payload.GetNewTokenBody{}
	isParseRequestPayloadSuccess := utils.ReadGinContextToPayload(c, &requestPayload)
	if !isParseRequestPayloadSuccess {
		return
	}

	if getNewTokenResult, getNewTokenError := keycloak.KeycloakGetNewToken(h.Ctx, requestPayload.Request.RefreshToken); getNewTokenError != nil {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			utils.ReturnResponse(
				c,
				constant.AuthenticateFailure,
				nil,
				getNewTokenError.Error(),
			),
		)
	} else {
		if getNewTokenResult.Error != "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				utils.ReturnResponse(
					c,
					constant.AuthenticateFailure,
					getNewTokenResult,
				),
			)
			return
		}
		protocolOpenidConnectTokenResponse := payload.ProtocolOpenidConnectTokenResponse(getNewTokenResult)
		c.JSON(
			http.StatusOK,
			utils.ReturnResponse(
				c,
				constant.Success,
				protocolOpenidConnectTokenResponse,
			),
		)
	}
}

func (h AuthenticateHandler) Logout(c *gin.Context) {

	ctx, isSuccess := utils.PrepareContext(c)
	if !isSuccess {
		return
	}
	h.Ctx = ctx

	requestPayload := payload.LogoutBody{}
	isParseRequestPayloadSuccess := utils.ReadGinContextToPayload(c, &requestPayload)
	if !isParseRequestPayloadSuccess {
		return
	}

	if revokeTokenResult, revokeTokenError := keycloak.KeycloakRevokeToken(h.Ctx, requestPayload.Request.RefreshToken); revokeTokenError != nil {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			utils.ReturnResponse(
				c,
				constant.AuthenticateFailure,
				nil,
				revokeTokenError.Error(),
			),
		)
	} else {
		revokeTokenErrorResponse := payload.KeycloakCommonErrorResponseResponse(revokeTokenResult)
		if revokeTokenErrorResponse.Error != "" {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				utils.ReturnResponse(
					c,
					constant.AuthenticateFailure,
					revokeTokenErrorResponse,
				),
			)
			return
		}
		c.JSON(
			http.StatusOK,
			utils.ReturnResponse(
				c,
				constant.Success,
				revokeTokenErrorResponse,
			),
		)
	}
}

func (h AuthenticateHandler) Register(c *gin.Context) {

	ctx, isSuccess := utils.PrepareContext(c, true)
	if !isSuccess {
		return
	}
	h.Ctx = ctx

	requestPayload := payload.RegisterRequestBody{}
	isParseRequestPayloadSuccess := utils.ReadGinContextToPayload(c, &requestPayload)
	if !isParseRequestPayloadSuccess {
		return
	}

	if registerUserError := keycloak.KeycloakUserRegister(h.Ctx, requestPayload.Request); registerUserError != nil {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			utils.ReturnResponse(
				c,
				constant.AuthenticateFailure,
				nil,
				registerUserError.Error(),
			),
		)
		return
	}
	c.JSON(
		http.StatusOK,
		utils.ReturnResponse(
			c,
			constant.Success,
			nil,
		),
	)
}

func (h AuthenticateHandler) ActiveAccount(c *gin.Context) {

	ctx, isSuccess := utils.PrepareContext(c, true)
	if !isSuccess {
		return
	}
	h.Ctx = ctx

	userId := c.Query("userId")
	username := c.Query("username")

	if userId == "" || username == "" {
		c.Data(
			http.StatusOK,
			constant.ContentTypeHTML,
			[]byte(`
				<h1>user id and username is empty</h1>
			`),
		)
		return
	}

	errorEnum := constant.Success

	transactionError := h.DB.WithContext(h.Ctx).Transaction(func(tx *gorm.DB) error {
		userAccountActivationLog := model.UserAccountActivationLog{}
		findUserAccountActivationLogResult := tx.Where(
			tx.Where(model.UserAccountActivationLog{
				UserId: userId,
			}),
		).Or(
			tx.Where(model.UserAccountActivationLog{
				Username: username,
			}),
		).Find(&userAccountActivationLog)
		if findUserAccountActivationLogResult.Error != nil {
			return findUserAccountActivationLogResult.Error
		}
		if userAccountActivationLog.BaseEntity.Id != 0 {
			errorEnum = constant.UserAccountAlreadyActived
			return errors.New(constant.UserAccountAlreadyActived.ErrorMessage)
		}

		if updateUserError := keycloak.KeycloakActiveAccount(h.Ctx, userId, username); updateUserError != nil {
			return updateUserError
		} else {
			userAccountActivationLog := model.UserAccountActivationLog{
				BaseEntity: utils.GenerateNewBaseEntity(h.Ctx),
				UserId:     userId,
				Username:   username,
			}
			tx.Save(&userAccountActivationLog)
		}
		return nil
	})
	if errorEnum.ErrorCode != 0 {
		c.Data(
			http.StatusOK,
			constant.ContentTypeHTML,
			[]byte(fmt.Sprintf(
				`
				<h1>Error: %s</h1>
			`, strconv.Itoa(errorEnum.ErrorCode)+" - "+errorEnum.ErrorMessage,
			)),
		)
		return
	}
	if transactionError != nil {
		c.Data(
			http.StatusOK,
			constant.ContentTypeHTML,
			[]byte(fmt.Sprintf(
				`
				<h1>Error: %s</h1>
			`, transactionError.Error(),
			)),
		)
		return
	}
	c.Data(
		http.StatusOK,
		constant.ContentTypeHTML,
		[]byte(
			`
			<h1>Active successfully</h1>
		`,
		),
	)
}
