package controller

import (
	"fast-storage-go-service/constant"
	"fast-storage-go-service/services/implement"
	"fast-storage-go-service/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func StorageController(router *gin.Engine, db *gorm.DB) {
	handler := &implement.StorageService{DB: db}

	storageRouter := router.Group(constant.BaseApiPath + "/storage")

	storageRouter.GET("/system_storage_status",
		utils.AuthenticationWithAuthorization([]string{}),
		utils.RequestLogger,
		utils.ResponseLogger,
		handler.SystemStorageStatus,
		utils.ErrorHandler)

}
