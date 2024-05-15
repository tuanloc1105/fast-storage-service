package implement

import (
	"context"
	"fast-storage-go-service/constant"
	"fast-storage-go-service/log"
	"fast-storage-go-service/payload"
	"fast-storage-go-service/utils"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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

func (h StorageHandler) GetAllElementInSpecificDirectory(c *gin.Context) {

	ctx, isSuccess := utils.PrepareContext(c)
	if !isSuccess {
		return
	}
	h.Ctx = ctx

	requestPayload := payload.GetAllElementInSpecificDirectoryBody{}
	isParseRequestPayloadSuccess := utils.ReadGinContextToPayload(c, &requestPayload)
	if !isParseRequestPayloadSuccess {
		return
	}

	systemRootFolder := log.GetSystemRootFolder()
	folderToView := handleProgressFolderToView(h.Ctx, systemRootFolder, requestPayload.Request.CurrentLocation)

	// check if use root folder is existing
	if _, directoryStatusError := os.Stat(folderToView); os.IsNotExist(directoryStatusError) {
		log.WithLevel(constant.Info, h.Ctx, "start to create folder %s", folderToView)
		makeDirectoryAllError := os.MkdirAll(folderToView, 0755)
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
	listFileInDirectoryCommand := fmt.Sprintf("ls -lh %s", folderToView)
	if listFileStdout, _, listFileError := utils.Shellout(h.Ctx, listFileInDirectoryCommand); listFileError != nil {
		if listFileError != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				utils.ReturnResponse(
					c,
					constant.CreateFolderError,
					nil,
					listFileError.Error(),
				),
			)
			return
		}
	} else {
		fmt.Println(listFileStdout)
		if strings.Contains(listFileStdout, "total 0") {
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
			if listFileStdout, _, listFileError := utils.Shellout(h.Ctx, listFileInDirectoryCommand); listFileError != nil {
				if listFileError != nil {
					c.AbortWithStatusJSON(
						http.StatusInternalServerError,
						utils.ReturnResponse(
							c,
							constant.ListFolderError,
							nil,
							listFileError.Error(),
						),
					)
					return
				}
			} else {
				var listOfFileInformation []payload.FileInformation
				commandResultLineArray := strings.Split(listFileStdout, "\n")
				for _, element := range commandResultLineArray {
					elementDetail := strings.Split(element, " !x&2 ")
					if len(elementDetail) != 3 {
						continue
					}
					elementType := "file"
					if strings.HasPrefix(elementDetail[0], "d") {
						elementType = "folder"
					}
					fileName, fileNameUnescapeError := url.QueryUnescape(elementDetail[2])
					if fileNameUnescapeError != nil {
						fileName = ""
					}
					fileInformation := payload.FileInformation{
						Size: elementDetail[1],
						Name: fileName,
						Type: elementType,
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

func (h StorageHandler) UploadFile(c *gin.Context) {

	ctx, isSuccess := utils.PrepareContext(c)
	if !isSuccess {
		return
	}
	h.Ctx = ctx

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

	systemRootFolder := log.GetSystemRootFolder()
	folderToView := handleProgressFolderToView(h.Ctx, systemRootFolder, folderLocation)

	fileUpload := multipartForm.File["file"]

	if fileUpload == nil || len(fileUpload) == 0 {
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

	file := fileUpload[0]

	fileUploadName := file.Filename

	fileUploadExtension := filepath.Ext(file.Filename)

	if fileUploadExtension != "" {
		fileUploadName = strings.Replace(fileUploadName, fileUploadExtension, "", -1)
	}

	// check if file is exist
	countNumberOfFileThatHaveTheSameNameCommand := fmt.Sprintf("ls -l %s | grep %s | wc -l", folderToView, fileUploadName+fileUploadExtension)
	countFileStdOut, _, countFileError := utils.Shellout(h.Ctx, countNumberOfFileThatHaveTheSameNameCommand)
	if countFileError != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.InternalFailure,
				nil,
				countFileError.Error(),
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
				numberOfFileIntConvertError.Error(),
			),
		)
		return
	}

	if numberOfFile > 0 {
		fileUploadName += " (" + uuid.New().String() + ")"
	}
	finalFileNameToSave := fileUploadName + fileUploadExtension
	finalFileLocation := folderToView + url.QueryEscape(finalFileNameToSave)

	if saveUploadedFileError := c.SaveUploadedFile(file, finalFileLocation); saveUploadedFileError != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.SaveFileError,
				nil,
				saveUploadedFileError.Error(),
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

func handleProgressFolderToView(ctx context.Context, systemRootFolder, inputCurrentLocation string) string {
	folderToView := ""

	if inputCurrentLocation == "" {
		folderToView = log.EnsureTrailingSlash(systemRootFolder + ctx.Value(constant.UserIdLogKey).(string))
	} else {
		var currentLocationFromRequestPayload string
		if strings.HasPrefix(inputCurrentLocation, "/") {
			currentLocationFromRequestPayload = inputCurrentLocation[1:]
		} else {
			currentLocationFromRequestPayload = inputCurrentLocation
		}
		folderToView = log.EnsureTrailingSlash(systemRootFolder+ctx.Value(constant.UserIdLogKey).(string)) + currentLocationFromRequestPayload
		folderToView = log.EnsureTrailingSlash(folderToView)
	}
	return folderToView
}
