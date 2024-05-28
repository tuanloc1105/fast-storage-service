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
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var FileAndFolderNameRegex *regexp.Regexp = regexp.MustCompile(`\d{2}:\d{2} (.+)$`)
var SizeOfFileInStatCommandResultRegex *regexp.Regexp = regexp.MustCompile(`Size: (\d+)`)

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
	var listOfFileInformation []payload.FileInformation = []payload.FileInformation{}

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
			for lineIndex, line := range lsCommandResultLineArray {
				if lineIndex < 1 {
					continue
				}
				fileName := ""
				fileSize := ""
				fileExtension := ""
				fileLastModifiedDate := ""
				fileType := "file"

				match := FileAndFolderNameRegex.FindStringSubmatch(line)
				if len(match) > 1 {
					fileName = match[1]
				}
				if fileName == "" {
					continue
				}
				statCommandForNameOrFolder := fmt.Sprintf("stat '%s'", fileName)
				if infomationOfNameOrFolderStdout, _, infomationOfNameOrFolderErrors := utils.ShelloutAtSpecificDirectory(h.Ctx, statCommandForNameOrFolder, folderToView, false, false); infomationOfNameOrFolderErrors != nil {
					log.WithLevel(constant.Warn, h.Ctx, "cannot execute stats: %s", infomationOfNameOrFolderErrors.Error())
					continue
				} else {
					listOfLineOfInformation := strings.Split(infomationOfNameOrFolderStdout, "\n")
					for informationLineIndex, informationLine := range listOfLineOfInformation {
						fmt.Println("line ", informationLineIndex, ": ", informationLine)
						// file/folder fileSize handler
						if informationLineIndex == 1 {
							if fileSizeMatch := SizeOfFileInStatCommandResultRegex.FindStringSubmatch(informationLine); fileSizeMatch != nil {
								if len(fileSizeMatch) > 1 {
									fileSizeInt64, fileSizeInt64ConvertError := strconv.ParseInt(fileSizeMatch[1], 10, 64)
									if fileSizeInt64ConvertError != nil {
										log.WithLevel(constant.Warn, h.Ctx, "cannot convert file size from string to int64 of file %s: %s", fileName, fileSizeInt64ConvertError.Error())
									}
									if BytesPerKB <= fileSizeInt64 && fileSizeInt64 < BytesPerMB {
										fileSize = fmt.Sprintf("%.4f %s", convertBytesToKB(fileSizeInt64), "KB")
									} else if BytesPerMB <= fileSizeInt64 && fileSizeInt64 < BytesPerGB {
										fileSize = fmt.Sprintf("%.4f %s", convertBytesToMB(fileSizeInt64), "MB")
									} else if BytesPerGB <= fileSizeInt64 {
										fileSize = fmt.Sprintf("%.4f %s", convertBytesToGB(fileSizeInt64), "GB")
									} else {
										fileSize = fmt.Sprintf("%d %s", fileSizeInt64, "byte(s)")
									}
									fmt.Println("fileSize is:", fileSize)
								}
							}
						}
						// file type handler
						if informationLineIndex == 3 {
							if strings.Contains(informationLine, "/d") {
								fileType = "folder"
							}
						}
						// last modified date handler
						if informationLineIndex == 5 {
							modifiedDateString := strings.TrimSpace(strings.Replace(informationLine, "Modify: ", "", -1))
							fmt.Println("modifiedDateString is:", modifiedDateString)
							if modifiedDate, modifiedDateParseError := time.Parse(constant.FileStatDateTimeLayout, modifiedDateString); modifiedDateParseError != nil {
								log.WithLevel(constant.Error, h.Ctx, "an error has been occurred while convert last modified time string: \n- %s", modifiedDateParseError.Error())
							} else {
								fileLastModifiedDate = modifiedDate.Format(constant.YyyyMmDdHhMmSsFormat)
							}

						}
					}
				}

				// file extension handler
				if fileType == "file" {
					getFileExtensionCommand := fmt.Sprintf("basename '%s' | awk -F. '{print $NF}'", fileName)
					fileExtensionStdOut, fileExtensionErrOut, fileExtensionError := utils.ShelloutAtSpecificDirectory(h.Ctx, getFileExtensionCommand, folderToView)
					if fileExtensionError != nil {
						log.WithLevel(constant.Error, h.Ctx, "an error has been occurred while geting file extension: \n- %s\n- %s", fileExtensionErrOut, fileExtensionError.Error())
					}
					fileExtension = strings.ToUpper(fileExtensionStdOut)
					fmt.Println("fileExtension is:", fileExtension)
				}

				fileInfo := payload.FileInformation{
					Size:             fileSize,
					Name:             fileName,
					Extension:        fileExtension,
					LastModifiedDate: fileLastModifiedDate,
					Type:             fileType,
				}
				listOfFileInformation = append(listOfFileInformation, fileInfo)
			}
		}
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
	createFileWithTouchCommand := "touch " + fileNameToCreate

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

	if strings.Contains(folderLocation, "..") || strings.Contains(folderLocation, ".") {
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
		finalFileLocation := folderToSaveFile + finalFileNameToSave

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

	folderLocation := c.Query("locationToDownload")
	credential := c.Query("credential")
	fileNameToDownloadFromRequest := c.Query("fileNameToDownload")

	if strings.Contains(folderLocation, "..") || strings.Contains(folderLocation, ".") {
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

	fileNameToDownload := fileNameToDownloadFromRequest
	finalFileName := fileNameToDownload
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
	// c.Data(http.StatusOK, constant.ContentTypeBinary, fileData)
	c.Writer.Write(fileData)
}

func (h StorageHandler) DownloadFolder(c *gin.Context) {

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

	folderLocation := c.Query("locationToDownload")
	credential := c.Query("credential")

	if strings.Contains(folderLocation, "..") || strings.Contains(folderLocation, ".") {
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

	if folderLocation == "" {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utils.ReturnResponse(
				c,
				constant.DataFormatError,
				nil,
				"`folderLocation` can not be empty",
			),
		)
		return
	}

	systemRootFolder := log.GetSystemRootFolder()
	folderToDownload := handleProgressFolderToView(h.Ctx, systemRootFolder, folderLocation)
	checkFolderCredentialError := handleCheckUserFolderSecurityActivities(h.Ctx, h.DB, folderToDownload, credential)
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

	// folderToDownload is the location that will be zip
	// but when zipping, the location to run command must be the outside folder
	// so this line of codes will get the outside folder location
	outsideFolderLocation := ""
	baseNameCommand := fmt.Sprintf("basename '%s' | awk -F. '{print $NF}'", folderToDownload)
	folderToBeZipped, _, baseNameCommandError := utils.Shellout(h.Ctx, baseNameCommand)
	if baseNameCommandError != nil {
		log.WithLevel(constant.Warn, h.Ctx, "can not get base name of folder with error %s", baseNameCommandError.Error())
	} else {
		outsideFolderLocation = strings.Replace(folderToDownload, folderToBeZipped, "", -1)
	}

	zipFolderCommand := fmt.Sprintf("zip -r '%s.zip' '%s/'", folderToBeZipped, folderToBeZipped)

	_, _, zipError := utils.ShelloutAtSpecificDirectory(h.Ctx, zipFolderCommand, outsideFolderLocation)
	if zipError != nil {
		log.WithLevel(constant.Warn, h.Ctx, "can not zip folder with error %s", zipError.Error())
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			utils.ReturnResponse(
				c,
				constant.ZipFolderError,
				nil,
			),
		)
		return
	}
	fileToReturnToClient, openFileError := os.Open(outsideFolderLocation + folderToBeZipped + ".zip")
	if openFileError != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.DownloadFileError,
				nil,
				"cannot open file "+outsideFolderLocation+folderToBeZipped+".zip"+" to download. "+openFileError.Error(),
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
				"cannot convert file "+outsideFolderLocation+folderToBeZipped+".zip"+" to download. "+readFileToReturnToClientError.Error(),
			),
		)
		return
	}

	c.Status(200)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", constant.ContentTypeBinary)
	c.Header("Content-Disposition", "attachment; filename="+folderToBeZipped+".zip")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Writer.Write(fileData)
	utils.ShelloutAtSpecificDirectory(h.Ctx, "rm -f "+folderToBeZipped+".zip", outsideFolderLocation)
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

	fileNameToDownload := requestPayload.Request.FileNameToRemove

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
