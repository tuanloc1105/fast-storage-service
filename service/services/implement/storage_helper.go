package implement

import (
	"context"
	"errors"
	"fast-storage-go-service/constant"
	"fast-storage-go-service/log"
	"fast-storage-go-service/model"
	"fast-storage-go-service/utils"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"gorm.io/gorm"
)

const KbPerMB = 1024.0
const MbPerGB = 1024.0

const BytesPerGB = 1024 * 1024 * 1024
const BytesPerMB = 1024 * 1024
const BytesPerKB = 1024

// handleProgressFolderToView will proccess the directory path from input,
//
// ensure that the result is the valid directory path
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

// handleCheckUserMaximumStorage check that if current use have been set a limitation of storage size
func handleCheckUserMaximumStorage(ctx context.Context, db *gorm.DB) error {
	if authorizedUsernameFromContext := ctx.Value(constant.UsernameLogKey); authorizedUsernameFromContext != nil {
		if currentUsername, isCurrentUsernameConvertableToString := authorizedUsernameFromContext.(string); isCurrentUsernameConvertableToString {
			userStorageMaximunSizeFromDb := model.UserStorageLimitationData{}
			db.WithContext(ctx).Where(
				model.UserStorageLimitationData{
					Username: currentUsername,
				},
			).Find(&userStorageMaximunSizeFromDb)
			if userStorageMaximunSizeFromDb.BaseEntity.Id != 0 {
				return nil
			} else {
				baseEntity := utils.GenerateNewBaseEntity(ctx)
				newUserStorageMaximunSizeToSaveToDb := model.UserStorageLimitationData{
					BaseEntity:         baseEntity,
					Username:           currentUsername,
					MaximunStorageSize: 1,
					StorageSizeUnit:    "GB",
				}
				saveDataResult := db.WithContext(ctx).Save(&newUserStorageMaximunSizeToSaveToDb)
				return saveDataResult.Error
			}
		}
		return errors.New("cannot convert current username data")
	}
	return errors.New("cannot determine current user")
}

func handleCheckUserMaximumStorageWhenUploading(ctx context.Context, db *gorm.DB, systemRootFolder string, fileUploadingSize int64) (float64, float64, error) {
	// if fileUploadingSize == 0 {
	// 	return 0, 0, nil
	// }
	if authorizedUsernameFromContext := ctx.Value(constant.UsernameLogKey); authorizedUsernameFromContext != nil {
		if currentUsername, isCurrentUsernameConvertableToString := authorizedUsernameFromContext.(string); isCurrentUsernameConvertableToString {
			userStorageMaximunSizeFromDb := model.UserStorageLimitationData{}
			db.WithContext(ctx).Where(
				model.UserStorageLimitationData{
					Username: currentUsername,
				},
			).Find(&userStorageMaximunSizeFromDb)
			if userStorageMaximunSizeFromDb.BaseEntity.Id == 0 {
				return 0, 0, errors.New("cannot check user storage limitation")
			}
			folderToCheckSize := systemRootFolder + ctx.Value(constant.UserIdLogKey).(string)
			checkFolderSizeCommand := fmt.Sprintf("du -s %s", folderToCheckSize)
			checkFolderSizeStdOut, _, checkFolderSizeError := utils.Shellout(ctx, checkFolderSizeCommand)
			if checkFolderSizeError != nil {
				return 0, 0, checkFolderSizeError
			}
			folderSizeInt64, convertFolderSizeToInt64Error := strconv.
				ParseInt(
					strings.Split(checkFolderSizeStdOut, "\t")[0],
					10,
					64,
				)
			if convertFolderSizeToInt64Error != nil {
				return 0, 0, convertFolderSizeToInt64Error
			}
			userMaximunMbStorage := convertGBToMB(float64(int64(userStorageMaximunSizeFromDb.MaximunStorageSize)))
			folderSizeMb := convertKBToMB(float64(folderSizeInt64))
			fileUploadingSizeMb := convertBytesToMB(fileUploadingSize)
			log.WithLevel(
				constant.Info,
				ctx,
				"user storage information when uploading file\n\t- user maximun storage: %.7f\n\t- current storage size of user: %.7f\n\t- file uploading size: %.7f",
				userMaximunMbStorage,
				folderSizeMb,
				fileUploadingSizeMb,
			)
			if folderSizeMb+fileUploadingSizeMb > userMaximunMbStorage {
				return 0, 0, errors.New("no more space to store file")
			}
			return userMaximunMbStorage, folderSizeMb, nil
		}
		return 0, 0, errors.New("cannot convert current username data")
	}
	return 0, 0, errors.New("cannot determine current user")
}

func handleCheckUserFolderSecurityActivities(ctx context.Context, db *gorm.DB, folderToCheck, credential string) error {

	currentTime := time.Now()
	secureFolderData := []model.UserFolderCredential{}
	folderSecureDataMatchWithInputFolder := model.UserFolderCredential{}

	db.WithContext(ctx).Where(
		model.UserFolderCredential{
			Username: ctx.Value(constant.UsernameLogKey).(string),
		},
	).Find(&secureFolderData)

	isInputFolderSecured := false

	for _, userFolderCredentialElement := range secureFolderData {
		if userFolderCredentialElement.Directory == folderToCheck || strings.Contains(folderToCheck, userFolderCredentialElement.Directory) {
			isInputFolderSecured = true
			folderSecureDataMatchWithInputFolder = userFolderCredentialElement
			break
		}
	}

	if !isInputFolderSecured {
		return nil
	}

	// check folder activity
	if currentTime.Sub(folderSecureDataMatchWithInputFolder.LastFolderActivitiesTime) < time.Duration(5)*time.Minute {
		folderSecureDataMatchWithInputFolder.LastFolderActivitiesTime = currentTime
		db.WithContext(ctx).Save(&folderSecureDataMatchWithInputFolder)
		return nil
	}

	// check folder credential
	var checkCredentialError error = nil
	if folderSecureDataMatchWithInputFolder.CredentialType == "OTP" {
		if _, handleOtpError := handleCheckUserOtp(ctx, db, credential); handleOtpError != nil {
			checkCredentialError = handleOtpError
		}
	} else {
		if comparePasswordError := utils.ComparePassword(credential, folderSecureDataMatchWithInputFolder.Credential); comparePasswordError != nil {
			checkCredentialError = comparePasswordError
		}
	}
	if checkCredentialError == nil {
		folderSecureDataMatchWithInputFolder.LastFolderActivitiesTime = currentTime
		db.WithContext(ctx).Save(&folderSecureDataMatchWithInputFolder)
	}
	return checkCredentialError
}

func folderIsSecure(ctx context.Context, db *gorm.DB, folderToCheck string) bool {

	secureFolderData := []model.UserFolderCredential{}

	db.WithContext(ctx).Where(
		model.UserFolderCredential{
			Username: ctx.Value(constant.UsernameLogKey).(string),
		},
	).Find(&secureFolderData)

	isInputFolderSecured := false

	for _, userFolderCredentialElement := range secureFolderData {
		if userFolderCredentialElement.Directory == folderToCheck || strings.Contains(folderToCheck, userFolderCredentialElement.Directory) {
			isInputFolderSecured = true
			break
		}
	}

	return isInputFolderSecured
}

func convertKBToMB(kb float64) float64 {
	return kb / KbPerMB
}

func convertGBToMB(gb float64) float64 {
	return gb * MbPerGB
}

func convertBytesToGB(bytes int64) float64 {
	return float64(bytes) / float64(BytesPerGB)
}

func convertBytesToMB(bytes int64) float64 {
	return float64(bytes) / float64(BytesPerMB)
}

func convertBytesToKB(bytes int64) float64 {
	return float64(bytes) / float64(BytesPerKB)
}

func removeNonAlpha(s string) string {
	result := []rune{}
	for _, char := range s {
		if unicode.IsLetter(char) || unicode.IsSpace(char) || unicode.IsNumber(char) {
			result = append(result, char)
		}
	}
	return string(result)
}
