package implement

import (
	"context"
	"fast-storage-go-service/config"
	"fast-storage-go-service/constant"
	"fast-storage-go-service/keycloak"
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
					listOfFileInformation,
				),
			)
			return
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
				fileEditable := false
				fileBirthDate := ""

				match := FileAndFolderNameRegex.FindStringSubmatch(line)
				if len(match) > 1 {
					fileName = match[1]
				}
				if fileName == "" {
					continue
				}

				// check if file is editable
				fileCommand := "file '" + fileName + "' "
				fileCommandStdout, _, fileCommandError := utils.ShelloutAtSpecificDirectory(h.Ctx, fileCommand, folderToView)
				if fileCommandError == nil {
					if strings.Contains(fileCommandStdout, "text") {
						fileEditable = true
					}
				}

				statCommandForNameOrFolder := fmt.Sprintf("stat '%s'", fileName)
				if infomationOfNameOrFolderStdout, _, infomationOfNameOrFolderErrors := utils.ShelloutAtSpecificDirectory(h.Ctx, statCommandForNameOrFolder, folderToView, true, true); infomationOfNameOrFolderErrors != nil {
					log.WithLevel(constant.Warn, h.Ctx, "cannot execute stats: %s", infomationOfNameOrFolderErrors.Error())
					continue
				} else {
					listOfLineOfInformation := strings.Split(infomationOfNameOrFolderStdout, "\n")
					if fileSizeMatch := SizeOfFileInStatCommandResultRegex.FindStringSubmatch(listOfLineOfInformation[1]); fileSizeMatch != nil {
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
							// fmt.Println("fileSize is:", fileSize)
						}
					}
					if strings.Contains(listOfLineOfInformation[3], "/d") {
						fileType = "folder"
					}
					modifiedDateString := strings.TrimSpace(strings.Replace(listOfLineOfInformation[5], "Modify: ", "", -1))
					// fmt.Println("modifiedDateString is:", modifiedDateString)
					if modifiedDate, modifiedDateParseError := time.Parse(constant.FileStatDateTimeLayout, modifiedDateString); modifiedDateParseError != nil {
						log.WithLevel(constant.Error, h.Ctx, "an error has been occurred while convert last modified time string: \n- %s", modifiedDateParseError.Error())
					} else {
						fileLastModifiedDate = modifiedDate.Format(constant.YyyyMmDdHhMmSsFormat)
					}

					birthDateString := strings.TrimSpace(strings.Replace(listOfLineOfInformation[7], "Birth: ", "", -1))
					// fmt.Println("birthDateString is:", birthDateString)
					if birthDate, birthDateParseError := time.Parse(constant.FileStatDateTimeLayout, birthDateString); birthDateParseError != nil {
						log.WithLevel(constant.Error, h.Ctx, "an error has been occurred while convert last modified time string: \n- %s", birthDateParseError.Error())
					} else {
						fileBirthDate = birthDate.Format(constant.YyyyMmDdHhMmSsFormat)
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
					// fmt.Println("fileExtension is:", fileExtension)
				}

				fileInfo := payload.FileInformation{
					Size:             fileSize,
					Name:             fileName,
					Extension:        fileExtension,
					LastModifiedDate: fileLastModifiedDate,
					Type:             fileType,
					Editable:         fileEditable,
					BirthDate:        fileBirthDate,
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
		if encryptionError := utils.FileEncryption(h.Ctx, finalFileLocation); encryptionError != nil {
			utils.Shellout(h.Ctx, fmt.Sprintln("rm", "-f", finalFileLocation))
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				utils.ReturnResponse(
					c,
					constant.FileCryptoError,
					nil,
					fileUploadName+fileUploadExtension+" has an error while uploading.",
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
	// finalFileName := fileNameToDownload
	if decryptionError := utils.FileDecryption(h.Ctx, folderToView+fileNameToDownload); decryptionError != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.FileCryptoError,
				nil,
				folderToView+fileNameToDownload+" has an error while downloading.",
			),
		)
		return
	}
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
		if encryptionError := utils.FileEncryption(h.Ctx, folderToView+fileNameToDownload); encryptionError != nil {
			log.WithLevel(
				constant.Error,
				h.Ctx,
				fmt.Sprintln(
					"cannot ecrypt file",
					folderToView+fileNameToDownload,
					". starting to remove. error: ",
					encryptionError,
				),
			)
			utils.Shellout(h.Ctx, fmt.Sprintln("rm", "-f", folderToView+fileNameToDownload))
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
	// c.Header("File-Name", finalFileName)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", constant.ContentTypeBinary)
	c.Header("Content-Disposition", "attachment; filename*=UTF-8''"+url.PathEscape(fileNameToDownload))
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
	archiveType := c.Query("archiveType")

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

	if decryptionError := utils.FileDecryption(h.Ctx, folderToDownload); decryptionError != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.FileCryptoError,
				nil,
				folderToDownload+" has an error while downloading.",
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
		outsideFolderLocation = strings.Replace(outsideFolderLocation, "//", "/", -1)
	}
	var zipFolderCommand string
	var zipExtension string
	zipFileName := folderToBeZipped + "-" + time.Now().Format(constant.RarFileTimeLayout)
	switch archiveType {
	case "zip":
		zipFolderCommand = fmt.Sprintf("zip -r '%s.zip' '%s/'", zipFileName, folderToBeZipped)
		zipExtension = ".zip"
	case "rar":
		zipFolderCommand = fmt.Sprintf("rar a -r -m5 '%s.rar' '%s/'", zipFileName, folderToBeZipped)
		zipExtension = ".rar"
	default:
		zipFolderCommand = fmt.Sprintf("zip -r '%s.zip' '%s/'", zipFileName, folderToBeZipped)
		zipExtension = ".zip"
	}

	_, _, zipError := utils.ShelloutAtSpecificDirectory(h.Ctx, zipFolderCommand, outsideFolderLocation, true, false)
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
	fileToReturnToClient, openFileError := os.Open(outsideFolderLocation + zipFileName + zipExtension)
	if openFileError != nil {
		utils.ShelloutAtSpecificDirectory(h.Ctx, "rm -f "+zipFileName+zipExtension, outsideFolderLocation)
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.DownloadFileError,
				nil,
				"cannot open file "+outsideFolderLocation+zipFileName+zipExtension+" to download. "+openFileError.Error(),
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
		utils.ShelloutAtSpecificDirectory(h.Ctx, "rm -f "+zipFileName+zipExtension, outsideFolderLocation)
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.DownloadFileError,
				nil,
				"cannot convert file "+outsideFolderLocation+zipFileName+zipExtension+" to download. "+readFileToReturnToClientError.Error(),
			),
		)
		return
	}

	c.Status(200)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", constant.ContentTypeBinary)
	c.Header("Content-Disposition", "attachment; filename="+zipFileName+zipExtension)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Writer.Write(fileData)
	utils.ShelloutAtSpecificDirectory(h.Ctx, "rm -f "+zipFileName+zipExtension, outsideFolderLocation)
	if encryptionError := utils.FileEncryption(h.Ctx, folderToDownload); encryptionError != nil {
		log.WithLevel(
			constant.Error,
			h.Ctx,
			fmt.Sprintln(
				"cannot ecrypt file",
				folderToDownload,
				". starting to remove. error: ",
				encryptionError,
			),
		)
		utils.Shellout(h.Ctx, fmt.Sprintln("rm", "-f", folderToDownload))
	}
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

	for _, fileNameToBeRemoved := range requestPayload.Request.FileNameToRemove {
		removeFileCommand := fmt.Sprintf("rm -rf '%s'", folderToView+fileNameToBeRemoved)

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

	if strings.Contains(requestPayload.Request.Folder, "\\") {
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

		// if password type is `OTP`, then check if user have already configured OTP
		if requestPayload.Request.CredentialType == "OTP" {
			if currentUserId, getCurrentUserIdError := utils.GetCurrentUserId(c); getCurrentUserIdError != nil {
				c.JSON(
					http.StatusForbidden,
					utils.ReturnResponse(
						c,
						constant.HashPasswordForSecuredFolderError,
						nil,
						getCurrentUserIdError.Error(),
					),
				)
				return
			} else {
				userOtpDataFoundInTheDb := model.UsersOtpData{}
				h.DB.WithContext(h.Ctx).Where(
					model.UsersOtpData{
						UserId: *currentUserId,
					},
				).Find(&userOtpDataFoundInTheDb)
				if userOtpDataFoundInTheDb.BaseEntity.Id == 0 {
					c.JSON(
						http.StatusForbidden,
						utils.ReturnResponse(
							c,
							constant.OtpError,
							nil,
							"User did not configure OTP",
						),
					)
					return
				}
			}
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

func (h StorageHandler) ReadTextFileContent(c *gin.Context) {
	ctx, isSuccess := utils.PrepareContext(c)
	if !isSuccess {
		return
	}
	h.Ctx = ctx

	folderLocation := c.Query("locationToRead")
	credential := c.Query("credential")
	fileNameToReadFromRequest := c.Query("fileNameToRead")

	if folderLocation == "" || fileNameToReadFromRequest == "" {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utils.ReturnResponse(
				c,
				constant.DataFormatError,
				nil,
				"`locationToRead` and `fileNameToRead` can not be empty",
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

	checkIfFileIsTextFileCommand := "file '" + fileNameToReadFromRequest + "'"
	checkIfFileIsTextStdout, _, checkIfFileIsTextError := utils.ShelloutAtSpecificDirectory(h.Ctx, checkIfFileIsTextFileCommand, folderToView)
	if checkIfFileIsTextError != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.InternalFailure,
				nil,
				"can not view file",
			),
		)
		return
	}

	if strings.Contains(checkIfFileIsTextStdout, "text") {
		catFileContentCommand := "cat " + fileNameToReadFromRequest
		catFileContentStdout, _, catFileContentError := utils.ShelloutAtSpecificDirectory(h.Ctx, catFileContentCommand, folderToView)
		if catFileContentError != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				utils.ReturnResponse(
					c,
					constant.InternalFailure,
					nil,
					"can not view file",
				),
			)
			return
		}

		// get file extension
		var extension string
		fileNameArrayAfterSplitDot := strings.Split(fileNameToReadFromRequest, ".")
		if len(fileNameArrayAfterSplitDot) == 0 {
			extension = "txt"
		} else {
			extension = fileNameArrayAfterSplitDot[len(fileNameArrayAfterSplitDot)-1]
		}

		fileContentInMarkdownSyntax := fmt.Sprintf("```%s\n%s\n```", extension, catFileContentStdout)

		c.Data(
			http.StatusOK,
			constant.ContentTypeText,
			[]byte(fileContentInMarkdownSyntax),
		)
	} else {
		c.Data(
			http.StatusOK,
			constant.ContentTypeText,
			[]byte(``),
		)
	}
}

func (h StorageHandler) EditTextFileContent(c *gin.Context) {
	ctx, isSuccess := utils.PrepareContext(c)
	if !isSuccess {
		return
	}
	h.Ctx = ctx

	if c.Request.Header.Get("Content-Type") != "text/plain" {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utils.ReturnResponse(
				c,
				constant.DataFormatError,
				nil,
				"Content-Type must be text/plain",
			),
		)
		return
	}

	folderLocation := c.Query("locationToEdit")
	credential := c.Query("credential")
	fileNameToEditFromRequest := c.Query("fileNameToEdit")
	if folderLocation == "" || fileNameToEditFromRequest == "" {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utils.ReturnResponse(
				c,
				constant.DataFormatError,
				nil,
				"`locationToEdit` and `fileNameToEdit` can not be empty",
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

	rawData, readRequestBodyError := io.ReadAll(c.Request.Body)
	if readRequestBodyError != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utils.ReturnResponse(
				c,
				constant.DataFormatError,
				nil,
				fmt.Sprintf("can not read content that you want to edit: %v", readRequestBodyError),
			),
		)
		return
	}

	// check if file exist or not
	checkFileExistenceCommand := "ls " + fileNameToEditFromRequest
	_, _, checkFileExistenceError := utils.ShelloutAtSpecificDirectory(h.Ctx, checkFileExistenceCommand, folderToView)
	if checkFileExistenceError != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.InternalFailure,
				nil,
				fmt.Sprintf("file does not exist: %v", checkFileExistenceError),
			),
		)
		return
	}

	contentToEdit := string(rawData)

	editFileCommand := fmt.Sprintf(
		`cat <<EOF > %s
%s
EOF`,
		fileNameToEditFromRequest,
		contentToEdit,
	)

	_, _, editFileError := utils.ShelloutAtSpecificDirectory(h.Ctx, editFileCommand, folderToView)
	if editFileError != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.InternalFailure,
				nil,
				fmt.Sprintf("can not edit file: %v", editFileError),
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

func (h StorageHandler) ShareFile(c *gin.Context) {
	ctx, isSuccess := utils.PrepareContext(c)
	if !isSuccess {
		return
	}
	h.Ctx = ctx
	requestPayload := payload.ShareFileBody{}
	isParseRequestPayloadSuccess := utils.ReadGinContextToPayload(c, &requestPayload)
	if !isParseRequestPayloadSuccess {
		return
	}
	if strings.Contains(requestPayload.Request.Folder, "..") ||
		strings.Contains(requestPayload.Request.File, "..") ||
		strings.Contains(requestPayload.Request.File, "/") {
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
	if requestPayload.Request.UserEmailToShare == nil || len(requestPayload.Request.UserEmailToShare) < 1 {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utils.ReturnResponse(
				c,
				constant.EmptyUserToShareFileOrFolderError,
				nil,
			),
		)
		return
	}

	listOfInvalidUsername := make(map[string]string)
	listOfUsernameToShare := []string{}

	// Get keycloak admin token

	adminProtocolOpenidConnectToken, protocolOpenidConnectTokenError := keycloak.KeycloakLogin(ctx, config.KeycloakAdminUsername, config.KeycloakAdminPassword)

	if protocolOpenidConnectTokenError != nil {
		log.WithLevel(constant.Error, h.Ctx, protocolOpenidConnectTokenError.Error())
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utils.ReturnResponse(
				c,
				constant.InternalFailure,
				nil,
				"An error has been occurred while check list of user to be shared",
			),
		)
		return
	}

	if adminProtocolOpenidConnectToken.Error != "" {
		log.WithLevel(constant.Error, h.Ctx, fmt.Sprint(adminProtocolOpenidConnectToken.Error, "-", adminProtocolOpenidConnectToken.ErrorDescription))
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utils.ReturnResponse(
				c,
				constant.InternalFailure,
				nil,
				"An error has been occurred while check list of user to be shared",
			),
		)
		return
	}
	if requestPayload.Request.UserEmailToShare[0] != "all" {
		for _, username := range requestPayload.Request.UserEmailToShare {
			searchResults, searchError := keycloak.KeycloakSearchUser(h.Ctx, adminProtocolOpenidConnectToken.AccessToken, username)
			if searchError != nil {
				listOfInvalidUsername[username] = searchError.Error()
				continue
			}
			if len(searchResults) < 1 {
				listOfInvalidUsername[username] = "user does not exist"
				continue
			}
			if len(searchResults) > 1 {
				listOfInvalidUsername[username] = "email is not unique"
				continue
			}
			listOfUsernameToShare = append(listOfUsernameToShare, searchResults[0].Username)
		}
	}
	if len(listOfInvalidUsername) > 0 {
		detailMessage := ""
		first := true
		for key, value := range listOfInvalidUsername {
			if !first {
				detailMessage += "<br>"
			}
			detailMessage += fmt.Sprintf("%s: %s", key, value)
			first = false
		}
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utils.ReturnResponse(
				c,
				constant.InvalidListOfUserEmailToShareError,
				nil,
				detailMessage,
			),
		)
		return
	}

	currentUsername, getCurrentUsernameError := utils.GetCurrentUsername(c)
	if getCurrentUsernameError != nil {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			utils.ReturnResponse(
				c,
				constant.Unauthorized,
				nil,
			),
		)
		return
	}

	listOfUserCanAccess := fmt.Sprint(
		",",
		strings.Join(listOfUsernameToShare, ","),
		",",
	)

	shareToken := utils.GenerateRandomString(10)

	systemRootFolder := log.GetSystemRootFolder()
	folderToView := handleProgressFolderToView(h.Ctx, systemRootFolder, requestPayload.Request.Folder)
	baseEntity := utils.GenerateNewBaseEntity(h.Ctx)
	currentTime := time.Now()
	currentTime = currentTime.Add(time.Duration(requestPayload.Request.TheTimeIntervalInMinutesToBeShared) * time.Minute)
	userFileAndFolderSharing := model.UserFileAndFolderSharing{
		BaseEntity:        baseEntity,
		Username:          *currentUsername,
		ListOfUsersShared: listOfUserCanAccess,
		Directory:         folderToView,
		FileName:          requestPayload.Request.File,
		ShareToken:        shareToken,
		ExpiredTime:       currentTime,
	}
	h.DB.WithContext(h.Ctx).Save(&userFileAndFolderSharing)

	c.JSON(
		http.StatusOK,
		utils.ReturnResponse(
			c,
			constant.Success,
			nil,
		),
	)
}

func (h StorageHandler) DownloadMultipleFile(c *gin.Context) {
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
	// archiveType := c.Query("archiveType")
	listOfFileToDownload := c.Query("listOfFileToDownload")

	if listOfFileToDownload == "" {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utils.ReturnResponse(
				c,
				constant.DataFormatError,
				nil,
				"`listOfFileToDownload` cannot be empty",
			),
		)
		return
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

	// check if there is a invalid file in `listOfFileToDownload`
	var listOfFullPathFileLocationToDownload []string = []string{}
	listOfFileToDownloadArray := strings.Split(listOfFileToDownload, "@231a")
	for _, listOfFileElement := range listOfFileToDownloadArray {
		currentFileElementTrim := strings.Trim(listOfFileElement, " ")
		fileFullPath := folderToDownload + currentFileElementTrim
		checkFileExistenceCommand := fmt.Sprint("ls -l ", "'", fileFullPath, "'")
		_, checkFileExistenceStderr, checkFileExistenceError := utils.Shellout(h.Ctx, checkFileExistenceCommand)
		if checkFileExistenceError != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				utils.ReturnResponse(
					c,
					constant.FileToDownloadInvalidError,
					nil,
					fmt.Sprint(checkFileExistenceError, " - ", checkFileExistenceStderr),
				),
			)
			return
		}
		// listOfFullPathFileLocationToDownload = append(listOfFullPathFileLocationToDownload, fileFullPath)
		if decryptionError := utils.FileDecryption(h.Ctx, currentFileElementTrim); decryptionError != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				utils.ReturnResponse(
					c,
					constant.FileCryptoError,
					nil,
					currentFileElementTrim+" has an error while downloading.",
				),
			)
			return
		}
		listOfFullPathFileLocationToDownload = append(listOfFullPathFileLocationToDownload, currentFileElementTrim)
	}
	log.WithLevel(constant.Debug, h.Ctx, fmt.Sprint("list of file to be downloaded\n", listOfFullPathFileLocationToDownload))
	zipFileName := fmt.Sprint(uuid.New().String(), ".zip")
	listOfFileCommand := ""
	for _, currentFilePath := range listOfFullPathFileLocationToDownload {
		listOfFileCommand += fmt.Sprint("'", currentFilePath, "'", " ")
	}
	zipListOfFileCommand := fmt.Sprintf("zip %s %s", zipFileName, strings.TrimSpace(listOfFileCommand))
	_, zipFileStderr, zipFileError := utils.ShelloutAtSpecificDirectory(h.Ctx, zipListOfFileCommand, folderToDownload)
	if zipFileError != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			utils.ReturnResponse(
				c,
				constant.ZipFolderError,
				nil,
				fmt.Sprint(zipFileError, " - ", zipFileStderr),
			),
		)
		return
	}
	fileToReturnToClient, openFileError := os.Open(folderToDownload + zipFileName)
	if openFileError != nil {
		utils.ShelloutAtSpecificDirectory(h.Ctx, "rm -f "+zipFileName, folderToDownload)
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.DownloadFileError,
				nil,
				"cannot open file "+zipFileName+" to download. "+openFileError.Error(),
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
		utils.ShelloutAtSpecificDirectory(h.Ctx, "rm -f "+zipFileName, folderToDownload)
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.DownloadFileError,
				nil,
				"cannot convert file "+zipFileName+" to download. "+readFileToReturnToClientError.Error(),
			),
		)
		return
	}
	c.Status(200)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", constant.ContentTypeBinary)
	c.Header("Content-Disposition", "attachment; filename="+zipFileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Writer.Write(fileData)
	utils.ShelloutAtSpecificDirectory(h.Ctx, "rm -f "+zipFileName, folderToDownload)
	for _, currentFile := range listOfFullPathFileLocationToDownload {
		if encryptionError := utils.FileEncryption(h.Ctx, currentFile); encryptionError != nil {
			log.WithLevel(
				constant.Error,
				h.Ctx,
				fmt.Sprintln(
					"cannot ecrypt file",
					currentFile,
					". starting to remove. error: ",
					encryptionError,
				),
			)
			utils.Shellout(h.Ctx, fmt.Sprintln("rm", "-f", currentFile))
		}
	}
}

func (h StorageHandler) CryptoEveryFolder(c *gin.Context) {
	ctx, isSuccess := utils.PrepareContext(c, true)
	if !isSuccess {
		return
	}
	apiKey := c.Request.Header.Get("api-key")
	encryptFolderApiKey := os.Getenv("ENCRYPT_FOLDER_API_KEY")
	if encryptFolderApiKey == "" {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.InternalFailure,
				nil,
			),
		)
		return
	}
	if apiKey == "" || strings.Compare(apiKey, encryptFolderApiKey) != 0 {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			utils.ReturnResponse(
				c,
				constant.Forbidden,
				nil,
				"Invalid api key to access this resource",
			),
		)
		return
	}
	h.Ctx = ctx
	requestPayload := payload.EncryptEveryFolderRequestBody{}
	isParseRequestPayloadSuccess := utils.ReadGinContextToPayload(c, &requestPayload)
	if !isParseRequestPayloadSuccess {
		return
	}
	rootFolderToEncryptOrDecrypt := log.GetSystemRootFolder()
	if requestPayload.Request.Encrypt {
		if encryptionError := utils.FileEncryption(h.Ctx, rootFolderToEncryptOrDecrypt); encryptionError != nil {
			log.WithLevel(
				constant.Error,
				h.Ctx,
				fmt.Sprintln(
					"cannot ecrypt file",
					rootFolderToEncryptOrDecrypt,
					". starting to remove. error: ",
					encryptionError,
				),
			)
		}
	} else {
		if decryptionError := utils.FileDecryption(h.Ctx, rootFolderToEncryptOrDecrypt); decryptionError != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				utils.ReturnResponse(
					c,
					constant.FileCryptoError,
					nil,
					rootFolderToEncryptOrDecrypt+" has an error while downloading.",
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
