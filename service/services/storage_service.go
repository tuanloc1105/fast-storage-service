package services

import "github.com/gin-gonic/gin"

type StorageService interface {
	SystemStorageStatus(c *gin.Context)
	GetAllElementInSpecificDirectory(c *gin.Context)
	UploadFile(c *gin.Context)
	DownloadFile(c *gin.Context)
	UserStorageStatus(c *gin.Context)
}
