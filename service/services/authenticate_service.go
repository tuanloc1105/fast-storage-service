package services

import (
	"fast-storage-go-service/constant"
	"fast-storage-go-service/keycloak"
	"fast-storage-go-service/payload"
	"fast-storage-go-service/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	KeycloakGrantType = "password"
	KeycloakScope     = "openid"
)

type AuthenticateHandler struct {
	DB *gorm.DB
}

func (h *AuthenticateHandler) Login(c *gin.Context) {

	ctx, isSuccess := utils.PrepareContext(c, true)

	if !isSuccess {
		return
	}

	requestPayload := payload.LoginRequestBody{}
	utils.ReadGinContextToPayload(c, &requestPayload)

	if loginResult, loginError := keycloak.KeycloakLogin(ctx, requestPayload.Request.Username, requestPayload.Request.Password); loginError != nil {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
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
		res := payload.ProtocolOpenidConnectTokenResponse{
			AccessToken:      loginResult.AccessToken,
			ExpiresIn:        loginResult.ExpiresIn,
			RefreshExpiresIn: loginResult.RefreshExpiresIn,
			RefreshToken:     loginResult.RefreshToken,
			TokenType:        loginResult.TokenType,
			IDToken:          loginResult.IDToken,
			NotBeforePolicy:  loginResult.NotBeforePolicy,
			SessionState:     loginResult.SessionState,
			Scope:            loginResult.Scope,
			Error:            loginResult.Error,
			ErrorDescription: loginResult.ErrorDescription,
		}
		c.JSON(
			http.StatusOK,
			utils.ReturnResponse(
				c,
				constant.Success,
				res,
			),
		)
	}
}

func (h *AuthenticateHandler) GetUserInfo(c *gin.Context) {

	ctx, isSuccess := utils.PrepareContext(c, true)

	if !isSuccess {
		return
	}

	if userInfoResult, userInfoError := keycloak.KeycloakGetUserInfo(ctx, c.GetHeader("Authorization")[7:]); userInfoError != nil {
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
		res := payload.OpenidConnectTokenIntrospectResponse{
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
				res,
			),
		)
	}
}
