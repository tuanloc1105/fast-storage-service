package implement

import (
	"context"
	"fast-storage-go-service/constant"
	"fast-storage-go-service/log"
	"fast-storage-go-service/model"
	"fast-storage-go-service/payload"
	"fast-storage-go-service/utils"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StorageHandler struct {
	DB  *gorm.DB
	Ctx context.Context
}

func (h StorageHandler) SystemStorageStatus(c *gin.Context) {

	ctx, isSuccess := utils.PrepareContext(c)
	if !isSuccess {
		return
	}
	h.Ctx = ctx

	nfsHost, isNfsHostSet := os.LookupEnv("NFS_HOST")
	if !isNfsHostSet {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.InternalFailure,
				"can not check storage status",
			),
		)
		return
	}

	dfCommand := "df -h | grep \"" + nfsHost + "\" | awk '{printf \"%s-%s-%s-%s\", $2, $3, $4, $5}'"

	shellStdout, _, shellError := utils.Shellout(h.Ctx, dfCommand)
	if shellError != nil {
		log.WithLevel(constant.Info, ctx, "an error has been occurred: %s", shellError.Error())
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.InternalFailure,
				shellError.Error(),
			),
		)
		return
	}

	statusArray := strings.Split(shellStdout, "-")

	if len(statusArray) < 4 {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.InternalFailure,
				nil,
				"not enought information",
			),
		)
		return
	}

	result := payload.SystemStorageStatus{
		Size:            statusArray[0],
		Used:            statusArray[1],
		Avail:           statusArray[2],
		UseInPercentage: statusArray[3],
	}

	c.JSON(
		http.StatusOK,
		utils.ReturnResponse(
			c,
			constant.Success,
			result,
		),
	)
}

func (h StorageHandler) UserStorageStatus(c *gin.Context) {

	ctx, isSuccess := utils.PrepareContext(c)
	if !isSuccess {
		return
	}
	h.Ctx = ctx
	systemRootFolder := log.GetSystemRootFolder()

	if maximunStorageSize, currentStorageSize, checkStorageSizeError := handleCheckUserMaximumStorageWhenUploading(h.Ctx, h.DB, systemRootFolder, 0); checkStorageSizeError != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.InternalFailure,
				checkStorageSizeError.Error(),
			),
		)
		return
	} else {
		result := payload.UserStorageStatus{
			MaximunSize: maximunStorageSize,
			Used:        currentStorageSize,
		}

		c.JSON(
			http.StatusOK,
			utils.ReturnResponse(
				c,
				constant.Success,
				result,
			),
		)
	}
}

func (h StorageHandler) GetAllElementInSpecificDirectory(c *gin.Context) {

	ctx, isSuccess := utils.PrepareContext(c)
	if !isSuccess {
		return
	}
	h.Ctx = ctx

	if checkMaximunStorageError := handleCheckUserMaximumStorage(h.Ctx, h.DB); checkMaximunStorageError != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.CheckMaximunStorageError,
				nil,
				checkMaximunStorageError.Error(),
			),
		)
		return
	}

	requestPayload := payload.GetAllElementInSpecificDirectoryBody{}
	isParseRequestPayloadSuccess := utils.ReadGinContextToPayload(c, &requestPayload)
	if !isParseRequestPayloadSuccess {
		return
	}

	if strings.Contains(requestPayload.Request.CurrentLocation, "..") {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			utils.ReturnResponse(
				c,
				constant.DataFormatError,
				nil,
				"Not accepted",
			),
		)
		return
	}

	systemRootFolder := log.GetSystemRootFolder()
	folderToView := handleProgressFolderToView(h.Ctx, systemRootFolder, requestPayload.Request.CurrentLocation)
	userRootFolder := handleProgressFolderToView(h.Ctx, systemRootFolder, constant.EmptyString)

	// check if use root folder is existing
	if _, directoryStatusError := os.Stat(userRootFolder); os.IsNotExist(directoryStatusError) {
		log.WithLevel(constant.Info, h.Ctx, "start to create folder %s", userRootFolder)
		makeDirectoryAllError := os.MkdirAll(userRootFolder, 0755)
		if makeDirectoryAllError != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				utils.ReturnResponse(
					c,
					constant.CreateFolderError,
					nil,
					makeDirectoryAllError.Error(),
				),
			)
			return
		}
	}

	checkFolderCredentialError := handleCheckUserFolderSecurityActivities(h.Ctx, h.DB, folderToView, requestPayload.Request.Credential)
	if checkFolderCredentialError != nil {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			utils.ReturnResponse(
				c,
				constant.SecureFolderInvalidCredentialError,
				nil,
				checkFolderCredentialError.Error(),
			),
		)
		return

	}

	listFileInDirectoryCommand := fmt.Sprintf("ls -lh '%s'", folderToView)
	if listFileStdout, _, listFileError := utils.Shellout(h.Ctx, listFileInDirectoryCommand); listFileError != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.FolderNotExistError,
				nil,
			),
		)
		return
	} else {
		lsCommandResultLineArray := strings.Split(listFileStdout, "\n")
		if strings.Contains(listFileStdout, "total 0") && len(lsCommandResultLineArray) <= 1 {
			c.JSON(
				http.StatusOK,
				utils.ReturnResponse(
					c,
					constant.Success,
					nil,
				),
			)
		} else {
			listFileInDirectoryCommand += " | awk 'NR>1{printf \"%s !x&2 %s !x&2 %s\\n\", $1, $5, $9}'"
			if listOfFileWithAwkStdout, _, listOfFileWithAwkError := utils.Shellout(h.Ctx, listFileInDirectoryCommand); listOfFileWithAwkError != nil {
				c.AbortWithStatusJSON(
					http.StatusInternalServerError,
					utils.ReturnResponse(
						c,
						constant.ListFolderError,
						nil,
						listOfFileWithAwkError.Error(),
					),
				)
				return
			} else {
				var listOfFileInformation []payload.FileInformation
				commandResultLineArray := strings.Split(listOfFileWithAwkStdout, "\n")
				for _, element := range commandResultLineArray {
					elementDetail := strings.Split(element, " !x&2 ")
					if len(elementDetail) != 3 {
						continue
					}
					elementType := "file"
					if strings.HasPrefix(elementDetail[0], "d") {
						elementType = "folder"
					}
					fileNameWithExtension, fileNameUnescapeError := url.QueryUnescape(elementDetail[2])
					if fileNameUnescapeError != nil {
						fileNameWithExtension = ""
					}
					fileSize := elementDetail[1]

					var fileInformation payload.FileInformation = payload.FileInformation{
						Size: fileSize,
						Name: fileNameWithExtension,
						Type: elementType,
					}

					statCommand := fmt.Sprintf("stat '%s' | grep \"Modify\"", elementDetail[2])
					statStdOut, statStdErr, statError := utils.ShelloutAtSpecificDirectory(h.Ctx, statCommand, folderToView)
					if statError != nil {
						log.WithLevel(constant.Error, h.Ctx, "an error has been occurred while get stat: \n- %s\n- %s", statStdErr, statError.Error())
					}
					if statStdOut != "" {
						modifiedDateString := strings.Replace(statStdOut, "Modify: ", "", -1)
						if modifiedDate, modifiedDateParseError := time.Parse(constant.FileStatDateTimeLayout, modifiedDateString); modifiedDateParseError != nil {
							log.WithLevel(constant.Error, h.Ctx, "an error has been occurred while convert last modified time string: \n- %s", modifiedDateParseError.Error())
						} else {
							fileInformation.LastModifiedDate = modifiedDate.Format(constant.YyyyMmDdHhMmSsFormat)
						}
					}

					if elementType == "file" {
						getFileExtensionCommand := fmt.Sprintf("basename '%s' | awk -F. '{print $NF}'", fileNameWithExtension)
						fileExtensionStdOut, fileExtensionErrOut, fileExtensionError := utils.ShelloutAtSpecificDirectory(h.Ctx, getFileExtensionCommand, folderToView)
						if fileExtensionError != nil {
							log.WithLevel(constant.Error, h.Ctx, "an error has been occurred while geting file extension: \n- %s\n- %s", fileExtensionErrOut, fileExtensionError.Error())
						}
						fileInformation.Extension = strings.ToUpper(fileExtensionStdOut)
					}

					listOfFileInformation = append(listOfFileInformation, fileInformation)
				}

				c.JSON(
					http.StatusOK,
					utils.ReturnResponse(
						c,
						constant.Success,
						listOfFileInformation,
					),
				)
			}
		}
	}
}

func (h StorageHandler) CreateFolder(c *gin.Context) {

	ctx, isSuccess := utils.PrepareContext(c)
	if !isSuccess {
		return
	}
	h.Ctx = ctx

	if checkMaximunStorageError := handleCheckUserMaximumStorage(h.Ctx, h.DB); checkMaximunStorageError != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.CheckMaximunStorageError,
				nil,
				checkMaximunStorageError.Error(),
			),
		)
		return
	}

	requestPayload := payload.CreateFolderBody{}
	isParseRequestPayloadSuccess := utils.ReadGinContextToPayload(c, &requestPayload)
	if !isParseRequestPayloadSuccess {
		return
	}
	if strings.Contains(requestPayload.Request.FolderToCreate, "..") {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			utils.ReturnResponse(
				c,
				constant.DataFormatError,
				nil,
				"Not accepted",
			),
		)
		return
	}

	systemRootFolder := log.GetSystemRootFolder()
	if strings.Contains(requestPayload.Request.FolderToCreate, "\\") {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utils.ReturnResponse(
				c,
				constant.DataFormatError,
				nil,
				"input cannot contain \\",
			),
		)
		return
	}
	folderToCreate := handleProgressFolderToView(h.Ctx, systemRootFolder, requestPayload.Request.FolderToCreate)

	// check if use root folder is existing
	if _, directoryStatusError := os.Stat(folderToCreate); os.IsNotExist(directoryStatusError) {
		log.WithLevel(constant.Info, h.Ctx, "start to create folder %s", folderToCreate)
		makeDirectoryAllError := os.MkdirAll(folderToCreate, 0755)
		if makeDirectoryAllError != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				utils.ReturnResponse(
					c,
					constant.CreateFolderError,
					nil,
					makeDirectoryAllError.Error(),
				),
			)
			return
		}
	}
	c.JSON(
		http.StatusOK,
		utils.ReturnResponse(
			c,
			constant.Success,
			nil,
		),
	)
}

func (h StorageHandler) RenameFileOrFolder(c *gin.Context) {

	ctx, isSuccess := utils.PrepareContext(c)
	if !isSuccess {
		return
	}
	h.Ctx = ctx

	if checkMaximunStorageError := handleCheckUserMaximumStorage(h.Ctx, h.DB); checkMaximunStorageError != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.CheckMaximunStorageError,
				nil,
				checkMaximunStorageError.Error(),
			),
		)
		return
	}

	requestPayload := payload.RenameFolderBody{}
	isParseRequestPayloadSuccess := utils.ReadGinContextToPayload(c, &requestPayload)
	if !isParseRequestPayloadSuccess {
		return
	}

	if strings.Contains(requestPayload.Request.OldFolderLocationName, "..") {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			utils.ReturnResponse(
				c,
				constant.DataFormatError,
				nil,
				"Not accepted",
			),
		)
		return
	}
	if strings.Contains(requestPayload.Request.NewFolderLocationName, "..") {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			utils.ReturnResponse(
				c,
				constant.DataFormatError,
				nil,
				"Not accepted",
			),
		)
		return
	}

	systemRootFolder := log.GetSystemRootFolder()
	if strings.Contains(requestPayload.Request.OldFolderLocationName, "\\") || strings.Contains(requestPayload.Request.NewFolderLocationName, "\\") {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utils.ReturnResponse(
				c,
				constant.DataFormatError,
				nil,
				"`oldFolderLocationName` or `newFolderLocationName` cannot contain \\",
			),
		)
		return
	}
	oldFolderName := handleProgressFolderToView(h.Ctx, systemRootFolder, requestPayload.Request.OldFolderLocationName)
	newFolderName := handleProgressFolderToView(h.Ctx, systemRootFolder, requestPayload.Request.NewFolderLocationName)

	// if the folder user want to be rename is a secure folder, prevent this behavior

	if folderIsSecure(h.Ctx, h.DB, oldFolderName) {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			utils.ReturnResponse(
				c,
				constant.RenameSecuredDirectoryError,
				nil,
			),
		)
		return
	}

	// check if folder or file is existing
	if _, directoryStatusError := os.Stat(oldFolderName); directoryStatusError != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.RenameNonexistentDirectoryError,
				nil,
				directoryStatusError.Error(),
			),
		)
		return
	}

	renameFolderCommand := fmt.Sprintf("mv %s %s", oldFolderName, newFolderName)

	_, renameFolderStdErr, renameFolderError := utils.Shellout(h.Ctx, renameFolderCommand)
	if renameFolderError != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.RenameFolderError,
				nil,
				fmt.Sprintf(
					"cannot rename folder\n  - error 1: %s\n  - error 2: %s",
					renameFolderStdErr,
					renameFolderError.Error(),
				),
			),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		utils.ReturnResponse(
			c,
			constant.Success,
			nil,
		),
	)
}

func (h StorageHandler) CreateFile(c *gin.Context) {

	ctx, isSuccess := utils.PrepareContext(c)
	if !isSuccess {
		return
	}
	h.Ctx = ctx

	if checkMaximunStorageError := handleCheckUserMaximumStorage(h.Ctx, h.DB); checkMaximunStorageError != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.CheckMaximunStorageError,
				nil,
				checkMaximunStorageError.Error(),
			),
		)
		return
	}

	requestPayload := payload.CreateFileBody{}
	isParseRequestPayloadSuccess := utils.ReadGinContextToPayload(c, &requestPayload)
	if !isParseRequestPayloadSuccess {
		return
	}
	if strings.Contains(requestPayload.Request.FolderToCreate, "..") {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			utils.ReturnResponse(
				c,
				constant.DataFormatError,
				nil,
				"Not accepted",
			),
		)
		return
	}

	systemRootFolder := log.GetSystemRootFolder()
	if strings.Contains(requestPayload.Request.FolderToCreate, "\\") {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utils.ReturnResponse(
				c,
				constant.DataFormatError,
				nil,
				"input cannot contain \\",
			),
		)
		return
	}
	folderToCreate := handleProgressFolderToView(h.Ctx, systemRootFolder, requestPayload.Request.FolderToCreate)

	// check if use root folder is existing
	if _, directoryStatusError := os.Stat(folderToCreate); os.IsNotExist(directoryStatusError) {
		log.WithLevel(constant.Info, h.Ctx, "start to create folder %s", folderToCreate)
		makeDirectoryAllError := os.MkdirAll(folderToCreate, 0755)
		if makeDirectoryAllError != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				utils.ReturnResponse(
					c,
					constant.CreateFolderError,
					nil,
					makeDirectoryAllError.Error(),
				),
			)
			return
		}
	}
	fileNameToCreate := requestPayload.Request.FileNameToCreate + "." + requestPayload.Request.FileExtension
	createFileWithTouchCommand := "touch " + url.QueryEscape(fileNameToCreate)

	_, createFileStdErr, createFileError := utils.ShelloutAtSpecificDirectory(h.Ctx, createFileWithTouchCommand, folderToCreate)
	if createFileStdErr != "" || createFileError != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.CreateFileError,
				nil,
				fmt.Sprintf(
					"cannot create file %s at %s",
					fileNameToCreate,
					folderToCreate,
				),
			),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		utils.ReturnResponse(
			c,
			constant.Success,
			nil,
		),
	)
}

func (h StorageHandler) UploadFile(c *gin.Context) {

	ctx, isSuccess := utils.PrepareContext(c)
	if !isSuccess {
		return
	}
	h.Ctx = ctx

	if checkMaximunStorageError := handleCheckUserMaximumStorage(h.Ctx, h.DB); checkMaximunStorageError != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.CheckMaximunStorageError,
				nil,
				checkMaximunStorageError.Error(),
			),
		)
		return
	}

	multipartForm, multipartFormError := c.MultipartForm()
	if multipartFormError != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.InternalFailure,
				nil,
				multipartFormError.Error(),
			),
		)
		return
	}

	folderLocation := ""
	folderLocationArray := multipartForm.Value["folderLocation"]

	if len(folderLocationArray) > 0 {
		folderLocation = folderLocationArray[0]
	}

	if strings.Contains(folderLocation, "..") {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			utils.ReturnResponse(
				c,
				constant.DataFormatError,
				nil,
				"Not accepted",
			),
		)
		return
	}

	credential := multipartForm.Value["credential"]

	systemRootFolder := log.GetSystemRootFolder()
	folderToSaveFile := handleProgressFolderToView(h.Ctx, systemRootFolder, folderLocation)
	var checkFolderCredentialError error = nil
	if len(credential) == 0 {
		checkFolderCredentialError = handleCheckUserFolderSecurityActivities(h.Ctx, h.DB, folderToSaveFile, "")
	} else {
		checkFolderCredentialError = handleCheckUserFolderSecurityActivities(h.Ctx, h.DB, folderToSaveFile, credential[0])
	}
	if checkFolderCredentialError != nil {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			utils.ReturnResponse(
				c,
				constant.SecureFolderInvalidCredentialError,
				nil,
				checkFolderCredentialError.Error(),
			),
		)
		return

	}

	fileUpload := multipartForm.File["file"]

	if fileUpload == nil || len(fileUpload) < 1 {
		c.AbortWithStatusJSON(
			http.StatusNotFound,
			utils.ReturnResponse(
				c,
				constant.EmptyFileInformationError,
				nil,
				"Empty file upload",
			),
		)
		return
	}

	for _, file := range fileUpload {
		fileUploadName := file.Filename

		if _, _, checkUploadingFileSize := handleCheckUserMaximumStorageWhenUploading(
			h.Ctx,
			h.DB,
			systemRootFolder,
			file.Size,
		); checkUploadingFileSize != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				utils.ReturnResponse(
					c,
					constant.UploadFileSizeExceeds,
					nil,
					checkUploadingFileSize.Error(),
				),
			)
			return
		}

		fileUploadExtension := filepath.Ext(file.Filename)

		if fileUploadExtension != "" {
			fileUploadName = strings.Replace(fileUploadName, fileUploadExtension, "", -1)
		}

		// check if file is exist
		countNumberOfFileThatHaveTheSameNameCommand := fmt.Sprintf("ls -l %s | grep '%s' | wc -l", folderToSaveFile, fileUploadName+fileUploadExtension)
		countFileStdOut, _, countFileError := utils.Shellout(h.Ctx, countNumberOfFileThatHaveTheSameNameCommand)
		if countFileError != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				utils.ReturnResponse(
					c,
					constant.InternalFailure,
					nil,
					fileUploadName+fileUploadExtension+" has an error while uploading. "+countFileError.Error(),
				),
			)
			return
		}
		numberOfFile, numberOfFileIntConvertError := strconv.Atoi(countFileStdOut)
		if numberOfFileIntConvertError != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				utils.ReturnResponse(
					c,
					constant.CountFileError,
					nil,
					fileUploadName+fileUploadExtension+" has an error while uploading. "+numberOfFileIntConvertError.Error(),
				),
			)
			return
		}

		if numberOfFile > 0 {
			fileUploadName += " (" + uuid.New().String() + ")"
		}
		finalFileNameToSave := fileUploadName + fileUploadExtension
		finalFileLocation := folderToSaveFile + url.QueryEscape(finalFileNameToSave)

		if saveUploadedFileError := c.SaveUploadedFile(file, finalFileLocation); saveUploadedFileError != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				utils.ReturnResponse(
					c,
					constant.SaveFileError,
					nil,
					fileUploadName+fileUploadExtension+" has an error while uploading. "+saveUploadedFileError.Error(),
				),
			)
			return
		}
	}

	c.JSON(
		http.StatusOK,
		utils.ReturnResponse(
			c,
			constant.Success,
			nil,
		),
	)
}

func (h StorageHandler) DownloadFile(c *gin.Context) {

	ctx, isSuccess := utils.PrepareContext(c)
	if !isSuccess {
		return
	}
	h.Ctx = ctx

	if checkMaximunStorageError := handleCheckUserMaximumStorage(h.Ctx, h.DB); checkMaximunStorageError != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.CheckMaximunStorageError,
				nil,
				checkMaximunStorageError.Error(),
			),
		)
		return
	}

	// requestPayload := payload.DownloadFileBody{}
	// isParseRequestPayloadSuccess := utils.ReadGinContextToPayload(c, &requestPayload)
	// if !isParseRequestPayloadSuccess {
	// 	return
	// }

	// locationToDownload
	// fileNameToDownload
	// credential

	folderLocation := c.Query("locationToDownload")
	credential := c.Query("credential")
	fileNameToDownloadFromRequest := c.Query("fileNameToDownload")

	if strings.Contains(folderLocation, "..") {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			utils.ReturnResponse(
				c,
				constant.DataFormatError,
				nil,
				"Not accepted",
			),
		)
		return
	}

	if folderLocation == "" || fileNameToDownloadFromRequest == "" {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utils.ReturnResponse(
				c,
				constant.DataFormatError,
				nil,
				"`folderLocation` and `fileNameToDownloadFromRequest` can not be empty",
			),
		)
		return
	}

	systemRootFolder := log.GetSystemRootFolder()
	folderToView := handleProgressFolderToView(h.Ctx, systemRootFolder, folderLocation)
	checkFolderCredentialError := handleCheckUserFolderSecurityActivities(h.Ctx, h.DB, folderToView, credential)
	if checkFolderCredentialError != nil {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			utils.ReturnResponse(
				c,
				constant.SecureFolderInvalidCredentialError,
				nil,
				checkFolderCredentialError.Error(),
			),
		)
		return

	}

	fileNameToDownload := url.QueryEscape(fileNameToDownloadFromRequest)
	finalFileName := ""
	if unescapedFileName, unescapedFileNameError := url.QueryUnescape(fileNameToDownload); unescapedFileNameError == nil {
		finalFileName = unescapedFileName
	} else {
		finalFileName = fileNameToDownload
	}
	// c.FileAttachment(folderToView+fileNameToDownload, fileNameToDownload)
	fileToReturnToClient, openFileError := os.Open(folderToView + fileNameToDownload)
	if openFileError != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.DownloadFileError,
				nil,
				"cannot open file "+folderToView+fileNameToDownload+" to download. "+openFileError.Error(),
			),
		)
		return
	}
	defer func(file *os.File) {
		closeFileError := file.Close()
		if closeFileError != nil {
			log.WithLevel(constant.Warn, h.Ctx, closeFileError.Error())
		}
	}(fileToReturnToClient)

	fileData, readFileToReturnToClientError := io.ReadAll(fileToReturnToClient)
	if readFileToReturnToClientError != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.DownloadFileError,
				nil,
				"cannot convert file "+folderToView+fileNameToDownload+" to download. "+readFileToReturnToClientError.Error(),
			),
		)
		return
	}

	c.Status(200)
	c.Header("File-Name", finalFileName)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", constant.ContentTypeBinary)
	c.Header("Content-Disposition", "attachment; filename="+fileNameToDownload)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Data(http.StatusOK, constant.ContentTypeBinary, fileData)
}

func (h StorageHandler) RemoveFile(c *gin.Context) {

	ctx, isSuccess := utils.PrepareContext(c)
	if !isSuccess {
		return
	}
	h.Ctx = ctx

	if checkMaximunStorageError := handleCheckUserMaximumStorage(h.Ctx, h.DB); checkMaximunStorageError != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.CheckMaximunStorageError,
				nil,
				checkMaximunStorageError.Error(),
			),
		)
		return
	}

	requestPayload := payload.RemoveFileBody{}
	isParseRequestPayloadSuccess := utils.ReadGinContextToPayload(c, &requestPayload)
	if !isParseRequestPayloadSuccess {
		return
	}
	if strings.Contains(requestPayload.Request.LocationToRemove, "..") {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			utils.ReturnResponse(
				c,
				constant.DataFormatError,
				nil,
				"Not accepted",
			),
		)
		return
	}

	if errorEnums, handleOtpError := handleCheckUserOtp(h.Ctx, h.DB, requestPayload.Request.OtpCredential); handleOtpError != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utils.ReturnResponse(
				c,
				errorEnums,
				nil,
				handleOtpError.Error(),
			),
		)
		return
	}

	folderLocation := requestPayload.Request.LocationToRemove

	systemRootFolder := log.GetSystemRootFolder()
	folderToView := handleProgressFolderToView(h.Ctx, systemRootFolder, folderLocation)
	checkFolderCredentialError := handleCheckUserFolderSecurityActivities(h.Ctx, h.DB, folderToView, requestPayload.Request.Credential)
	if checkFolderCredentialError != nil {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			utils.ReturnResponse(
				c,
				constant.SecureFolderInvalidCredentialError,
				nil,
				checkFolderCredentialError.Error(),
			),
		)
		return

	}

	fileNameToDownload := url.QueryEscape(requestPayload.Request.FileNameToRemove)

	removeFileCommand := fmt.Sprintf("rm -rf %s", folderToView+fileNameToDownload)

	if _, _, removeFileCommandError := utils.Shellout(h.Ctx, removeFileCommand); removeFileCommandError != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.RemoveFileError,
				nil,
				removeFileCommandError.Error(),
			),
		)
		return
	} else {
		c.JSON(
			http.StatusOK,
			utils.ReturnResponse(
				c,
				constant.Success,
				nil,
			),
		)
	}
}

func (h StorageHandler) SetPasswordForFolder(c *gin.Context) {

	ctx, isSuccess := utils.PrepareContext(c)
	if !isSuccess {
		return
	}
	h.Ctx = ctx

	if checkMaximunStorageError := handleCheckUserMaximumStorage(h.Ctx, h.DB); checkMaximunStorageError != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.CheckMaximunStorageError,
				nil,
				checkMaximunStorageError.Error(),
			),
		)
		return
	}

	requestPayload := payload.SetPasswordForFolderBody{}
	isParseRequestPayloadSuccess := utils.ReadGinContextToPayload(c, &requestPayload)
	if !isParseRequestPayloadSuccess {
		return
	}
	if strings.Contains(requestPayload.Request.Folder, "..") {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			utils.ReturnResponse(
				c,
				constant.DataFormatError,
				nil,
				"Not accepted",
			),
		)
		return
	}

	systemRootFolder := log.GetSystemRootFolder()
	folderToSecure := handleProgressFolderToView(h.Ctx, systemRootFolder, requestPayload.Request.Folder)

	if requestPayload.Request.CredentialType != "OTP" && requestPayload.Request.CredentialType != "PASSWORD" {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.DataFormatError,
				nil,
				"`credentialType` is invalid. only accept `OTP` or `PASSWORD`",
			),
		)
		return
	}

	if requestPayload.Request.CredentialType == "PASSWORD" && requestPayload.Request.Credential == "" {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.DataFormatError,
				nil,
				"`credentialType` is `PASSWORD`. a credential must be sent in input",
			),
		)
		return
	}

	if strings.Contains(requestPayload.Request.Folder, "/") || strings.Contains(requestPayload.Request.Folder, "\\") {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utils.ReturnResponse(
				c,
				constant.DataFormatError,
				nil,
				"input cannot contain / or \\",
			),
		)
		return
	}

	if _, directoryStatusError := os.Stat(folderToSecure); os.IsNotExist(directoryStatusError) {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.FolderNotExistError,
				nil,
				fmt.Sprintf("`%s` does not exist. %s", folderToSecure, directoryStatusError.Error()),
			),
		)
		return
	}

	userFolderPasswordInDatabase := model.UserFolderCredential{}
	h.DB.WithContext(h.Ctx).Where(model.UserFolderCredential{
		Username:  h.Ctx.Value(constant.UsernameLogKey).(string),
		Directory: folderToSecure,
	},
	).
		Find(&userFolderPasswordInDatabase)

	if userFolderPasswordInDatabase.BaseEntity.Id != 0 {
		c.JSON(
			http.StatusForbidden,
			utils.ReturnResponse(
				c,
				constant.FolderAlreadySecureError,
				nil,
			),
		)
		return
	}

	if encryptPasswordError := utils.EncryptPasswordPointer(&requestPayload.Request.Credential); encryptPasswordError != nil {
		c.JSON(
			http.StatusForbidden,
			utils.ReturnResponse(
				c,
				constant.HashPasswordForSecuredFolderError,
				nil,
			),
		)
		return
	} else {
		baseEntity := utils.GenerateNewBaseEntity(h.Ctx)
		userPasswordCredential := model.UserFolderCredential{
			BaseEntity:               baseEntity,
			Username:                 h.Ctx.Value(constant.UsernameLogKey).(string),
			Directory:                folderToSecure,
			Credential:               requestPayload.Request.Credential,
			CredentialType:           requestPayload.Request.CredentialType,
			LastFolderActivitiesTime: baseEntity.CreatedAt,
		}

		h.DB.WithContext(h.Ctx).Save(&userPasswordCredential)

		c.JSON(
			http.StatusOK,
			utils.ReturnResponse(
				c,
				constant.Success,
				nil,
			),
		)
	}
}

func (h StorageHandler) CheckSecureFolderStatus(c *gin.Context) {

	ctx, isSuccess := utils.PrepareContext(c)
	if !isSuccess {
		return
	}
	h.Ctx = ctx

	if checkMaximunStorageError := handleCheckUserMaximumStorage(h.Ctx, h.DB); checkMaximunStorageError != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.CheckMaximunStorageError,
				nil,
				checkMaximunStorageError.Error(),
			),
		)
		return
	}

	requestPayload := payload.CheckSecureFolderStatusBody{}
	isParseRequestPayloadSuccess := utils.ReadGinContextToPayload(c, &requestPayload)
	if !isParseRequestPayloadSuccess {
		return
	}
	if strings.Contains(requestPayload.Request.Folder, "..") {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			utils.ReturnResponse(
				c,
				constant.DataFormatError,
				nil,
				"Not accepted",
			),
		)
		return
	}
	systemRootFolder := log.GetSystemRootFolder()
	folderToSecure := handleProgressFolderToView(h.Ctx, systemRootFolder, requestPayload.Request.Folder)

	currentTime := time.Now()

	secureFolderData := []model.UserFolderCredential{}
	folderSecureDataMatchWithInputFolder := model.UserFolderCredential{}

	h.DB.WithContext(h.Ctx).Where(
		model.UserFolderCredential{
			Username: ctx.Value(constant.UsernameLogKey).(string),
		},
	).Find(&secureFolderData)

	for _, userFolderCredentialElement := range secureFolderData {
		if userFolderCredentialElement.Directory == folderToSecure || strings.Contains(folderToSecure, userFolderCredentialElement.Directory) {
			folderSecureDataMatchWithInputFolder = userFolderCredentialElement
			break
		}
	}

	if folderSecureDataMatchWithInputFolder.BaseEntity.Id == 0 || currentTime.Sub(folderSecureDataMatchWithInputFolder.LastFolderActivitiesTime) < time.Duration(5)*time.Minute {
		c.JSON(
			http.StatusOK,
			utils.ReturnResponse(
				c,
				constant.Success,
				true,
			),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		utils.ReturnResponse(
			c,
			constant.Success,
			false,
		),
	)
}
