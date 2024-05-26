package services

import "github.com/gin-gonic/gin"

type StorageService interface {
	SystemStorageStatus(c *gin.Context)
	GetAllElementInSpecificDirectory(c *gin.Context)
	UploadFile(c *gin.Context)
	DownloadFile(c *gin.Context)
	UserStorageStatus(c *gin.Context)
	RemoveFile(c *gin.Context)
	CreateFolder(c *gin.Context)
	SetPasswordForFolder(c *gin.Context)
	CheckSecureFolderStatus(c *gin.Context)
	CreateFile(c *gin.Context)
	RenameFileOrFolder(c *gin.Context)
}
