package services

import "github.com/gin-gonic/gin"

type StorageService interface {
	SystemStorageStatus(c *gin.Context)
	GetAllElementInSpecificDirectory(c *gin.Context)
}
