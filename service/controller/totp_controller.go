package controller

import (
	"fast-storage-go-service/constant"
	"fast-storage-go-service/services"
	"fast-storage-go-service/services/implement"
	"fast-storage-go-service/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TotpController(router *gin.Engine, db *gorm.DB) {
	var handler services.TotpService = &implement.TotpHandler{DB: db}

	totpRouter := router.Group(constant.BaseApiPath + "/otp")

	totpRouter.GET("/generate_qr_code",
		utils.AuthenticationWithAuthorization([]string{}),
		utils.RequestLogger,
		utils.ResponseLogger,
		handler.GenerateQrCode,
		utils.ErrorHandler)

}
