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
}