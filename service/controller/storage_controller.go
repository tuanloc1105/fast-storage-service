package controller

import (
	"fast-storage-go-service/constant"
	"fast-storage-go-service/services"
	"fast-storage-go-service/services/implement"
	"fast-storage-go-service/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func StorageController(router *gin.Engine, db *gorm.DB) {
	var handler services.StorageService = &implement.StorageHandler{DB: db}

	storageRouter := router.Group(constant.BaseApiPath + "/storage")

	storageRouter.GET("/system_storage_status",
		utils.AuthenticationWithAuthorization([]string{}),
		utils.RequestLogger,
		utils.ResponseLogger,
		handler.SystemStorageStatus,
		utils.ErrorHandler)

	storageRouter.GET("/user_storage_status",
		utils.AuthenticationWithAuthorization([]string{}),
		utils.RequestLogger,
		utils.ResponseLogger,
		handler.UserStorageStatus,
		utils.ErrorHandler)

	storageRouter.POST("/get_all_element_in_specific_directory",
		utils.AuthenticationWithAuthorization([]string{}),
		utils.RequestLogger,
		utils.ResponseLogger,
		handler.GetAllElementInSpecificDirectory,
		utils.ErrorHandler)

	storageRouter.POST("/upload_file",
		utils.AuthenticationWithAuthorization([]string{}),
		utils.RequestLogger,
		utils.ResponseLogger,
		handler.UploadFile,
		utils.ErrorHandler)

	storageRouter.GET("/download_file",
		utils.GetTokenInParamAndSetToHeader(),
		utils.AuthenticationWithAuthorization([]string{}),
		utils.RequestLogger,
		utils.ResponseLogger,
		handler.DownloadFile,
		utils.ErrorHandler)

	storageRouter.GET("/download_folder",
		utils.GetTokenInParamAndSetToHeader(),
		utils.AuthenticationWithAuthorization([]string{}),
		utils.RequestLogger,
		utils.ResponseLogger,
		handler.DownloadFolder,
		utils.ErrorHandler)

	storageRouter.POST("/remove_file",
		utils.AuthenticationWithAuthorization([]string{}),
		utils.RequestLogger,
		utils.ResponseLogger,
		handler.RemoveFile,
		utils.ErrorHandler)

	storageRouter.POST("/create_folder",
		utils.AuthenticationWithAuthorization([]string{}),
		utils.RequestLogger,
		utils.ResponseLogger,
		handler.CreateFolder,
		utils.ErrorHandler)

	storageRouter.POST("/rename_file_or_folder",
		utils.AuthenticationWithAuthorization([]string{}),
		utils.RequestLogger,
		utils.ResponseLogger,
		handler.RenameFileOrFolder,
		utils.ErrorHandler)

	storageRouter.POST("/create_file",
		utils.AuthenticationWithAuthorization([]string{}),
		utils.RequestLogger,
		utils.ResponseLogger,
		handler.CreateFile,
		utils.ErrorHandler)

	storageRouter.POST("/set_password_for_folder",
		utils.AuthenticationWithAuthorization([]string{}),
		utils.RequestLogger,
		utils.ResponseLogger,
		handler.SetPasswordForFolder,
		utils.ErrorHandler)

	storageRouter.POST("/check_secure_folder_status",
		utils.AuthenticationWithAuthorization([]string{}),
		utils.RequestLogger,
		utils.ResponseLogger,
		handler.CheckSecureFolderStatus,
		utils.ErrorHandler)

	storageRouter.GET("/read_text_file_content",
		utils.AuthenticationWithAuthorization([]string{}),
		utils.RequestLogger,
		utils.ResponseLogger,
		handler.ReadTextFileContent,
		utils.ErrorHandler)

	storageRouter.POST("/edit_text_file_content",
		utils.AuthenticationWithAuthorization([]string{}),
		utils.RequestLogger,
		utils.ResponseLogger,
		handler.EditTextFileContent,
		utils.ErrorHandler)

	storageRouter.POST("/share_file",
		utils.AuthenticationWithAuthorization([]string{}),
		utils.RequestLogger,
		utils.ResponseLogger,
		handler.ShareFile,
		utils.ErrorHandler)

	storageRouter.GET("/download_multiple_file",
		utils.GetTokenInParamAndSetToHeader(),
		utils.AuthenticationWithAuthorization([]string{}),
		utils.RequestLogger,
		utils.ResponseLogger,
		handler.DownloadMultipleFile,
		utils.ErrorHandler)

	storageRouter.POST("/crypto_every_folder",
		utils.RequestLogger,
		utils.ResponseLogger,
		handler.CryptoEveryFolder,
		utils.ErrorHandler)

	storageRouter.POST("/search_file",
		utils.AuthenticationWithAuthorization([]string{}),
		utils.RequestLogger,
		utils.ResponseLogger,
		handler.SearchFile,
		utils.ErrorHandler)

	storageRouter.POST("/read_image_file",
		utils.AuthenticationWithAuthorization([]string{}),
		utils.RequestLogger,
		utils.ResponseLogger,
		handler.ReadImageFile,
		utils.ErrorHandler)

	storageRouter.POST("/cut_or_copy",
		utils.AuthenticationWithAuthorization([]string{}),
		utils.RequestLogger,
		utils.ResponseLogger,
		handler.CutOrCopy,
		utils.ErrorHandler)
}
