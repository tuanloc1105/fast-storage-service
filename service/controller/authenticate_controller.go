package controller

import (
	"fast-storage-go-service/constant"
	"fast-storage-go-service/services"
	"fast-storage-go-service/services/implement"
	"fast-storage-go-service/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthenticateController(router *gin.Engine, db *gorm.DB) {
	var handler services.AuthenticateService = &implement.AuthenticateHandler{DB: db}

	authenticateRouter := router.Group(constant.BaseApiPath + "/auth")

	authenticateRouter.POST("/login",
		utils.RequestLogger,
		utils.ResponseLogger,
		handler.Login,
		utils.ErrorHandler)

	authenticateRouter.GET("/get_user_info",
		utils.AuthenticationWithAuthorization([]string{}),
		utils.RequestLogger,
		utils.ResponseLogger,
		handler.GetUserInfo,
		utils.ErrorHandler)

	authenticateRouter.POST("/get_new_token",
		// utils.AuthenticationWithAuthorization([]string{}),
		utils.RequestLogger,
		utils.ResponseLogger,
		handler.GetNewToken,
		utils.ErrorHandler)

	authenticateRouter.POST("/logout",
		utils.AuthenticationWithAuthorization([]string{}),
		utils.RequestLogger,
		utils.ResponseLogger,
		handler.Logout,
		utils.ErrorHandler)

	authenticateRouter.POST("/register",
		utils.RequestLogger,
		utils.ResponseLogger,
		handler.Register,
		utils.ErrorHandler)

	authenticateRouter.GET("/active_account",
		utils.RequestLogger,
		utils.ResponseLogger,
		handler.ActiveAccount,
		utils.ErrorHandler)

}
