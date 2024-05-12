package implement

import (
	"context"
	"fast-storage-go-service/constant"
	"fast-storage-go-service/model"
	"fast-storage-go-service/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TotpHandler struct {
	DB  *gorm.DB
	Ctx context.Context
}

func generateSecret(ctx context.Context) (string, error) {
	if generateSecretShellStd, _, generateSecretError := utils.Shellout(ctx, "java -jar additional_source_code/two-factor-auth.jar \"GENERATE_BASE32_SECRET\""); generateSecretError != nil {
		return "", generateSecretError
	} else {
		return generateSecretShellStd, nil
	}
}

func (h TotpHandler) GenerateQrCode(c *gin.Context) {

	ctx, isSuccess := utils.PrepareContext(c)
	if !isSuccess {
		return
	}
	h.Ctx = ctx

	// check if user already have configured OTP
	userOtpDataFoundInDatabase := model.UsersOtpData{}
	h.DB.WithContext(h.Ctx).
		Where(
			model.UsersOtpData{
				UserId: h.Ctx.Value(constant.UserIdLogKey).(string),
				BaseEntity: model.BaseEntity{
					Active: utils.GetPointerOfAnyValue(true),
				},
			},
		).
		Find(&userOtpDataFoundInDatabase)

	if userOtpDataFoundInDatabase.BaseEntity.Id != 0 {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			utils.ReturnResponse(
				c,
				constant.UserAlreadyConfigureOtp,
				nil,
			),
		)
		return
	}

	if secretKey, secretKeyGeneratorError := generateSecret(ctx); secretKeyGeneratorError != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			utils.ReturnResponse(
				c,
				constant.InternalFailure,
				secretKeyGeneratorError.Error(),
			),
		)
		return
	} else {
		qrCodeLabel := "fs-service-" + uuid.New().String()
		generateQrCodeCommand := fmt.Sprintf(
			"java -jar additional_source_code/two-factor-auth.jar \"GENERATE_QR_IMAGE_URL\" \"%s\" \"%s\"",
			secretKey,
			qrCodeLabel,
		)
		if qrcodeShellOut, _, qrcodeError := utils.Shellout(h.Ctx, generateQrCodeCommand); qrcodeError != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				utils.ReturnResponse(
					c,
					constant.InternalFailure,
					qrcodeError.Error(),
				),
			)
			return
		} else {
			qrCodeImageBase64Data := strings.Replace(qrcodeShellOut, "\n", "", -1)
			baseEntity := utils.GenerateNewBaseEntity(h.Ctx)
			userOtpData := model.UsersOtpData{
				BaseEntity:                   baseEntity,
				UserId:                       h.Ctx.Value(constant.UserIdLogKey).(string),
				UserOtpSecretData:            secretKey,
				UserOtpQrCodeImageBase64Data: qrCodeImageBase64Data,
			}
			h.DB.WithContext(h.Ctx).Transaction(func(tx *gorm.DB) error {
				saveUserOtpDataResult := tx.Save(&userOtpData)
				if saveUserOtpDataResult.Error != nil {
					return saveUserOtpDataResult.Error
				}
				return nil
			})
			c.JSON(
				http.StatusOK,
				utils.ReturnResponse(
					c,
					constant.Success,
					qrCodeImageBase64Data,
				),
			)
		}
	}

}

func (h TotpHandler) GenerateTotp(c *gin.Context) {

	ctx, isSuccess := utils.PrepareContext(c)
	if !isSuccess {
		return
	}
	h.Ctx = ctx
	utils.Shellout(h.Ctx, "java -jar additional_source_code/two-factor-auth.jar \"GENERATE_CURRENT_OTP\" \"VOFZHNG45ATPCO4K\"")
}
