package services

import "github.com/gin-gonic/gin"

type StorageService interface {
	SystemStorageStatus(c *gin.Context)
	UserStorageStatus(c *gin.Context)
	GetAllElementInSpecificDirectory(c *gin.Context)
	CreateFolder(c *gin.Context)
	RenameFileOrFolder(c *gin.Context)
	CreateFile(c *gin.Context)
	UploadFile(c *gin.Context)
	DownloadFile(c *gin.Context)
	DownloadFolder(c *gin.Context)
	RemoveFile(c *gin.Context)
	SetPasswordForFolder(c *gin.Context)
	CheckSecureFolderStatus(c *gin.Context)
	ReadTextFileContent(c *gin.Context)
	EditTextFileContent(c *gin.Context)
	ShareFile(c *gin.Context)
}
