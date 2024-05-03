package controller

import (
	"fast-storage-go-service/constant"
	"fast-storage-go-service/services"
	"fast-storage-go-service/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthenticateController(router *gin.Engine, db *gorm.DB) {
	handler := &services.AuthenticateHandler{DB: db}

	authenticateRouter := router.Group(constant.BaseApiPath + "/auth")
	authenticateRouter.POST("/login", utils.RequestLogger, utils.ResponseLogger, handler.Login, utils.ErrorHandler)
	authenticateRouter.GET("/get_user_info", utils.AuthenticationWithAuthorization([]string{}), utils.RequestLogger, utils.ResponseLogger, handler.GetUserInfo, utils.ErrorHandler)
	authenticateRouter.POST("/get_new_token", utils.RequestLogger, utils.ResponseLogger, handler.GetNewToken, utils.ErrorHandler)
}
