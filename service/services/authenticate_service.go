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
		c.JSON(
			http.StatusOK,
			utils.ReturnResponse(
				c,
				constant.Success,
				loginResult,
			),
		)
	}
}
