package implement

import (
	"context"
	"errors"
	"fast-storage-go-service/constant"
	"fast-storage-go-service/log"
	"fast-storage-go-service/model"

	"gorm.io/gorm"
)

func handleCheckUserOtp(ctx context.Context, db *gorm.DB, inputOtp string) (constant.ErrorEnums, error) {
	// check if user is enable OTP
	userOtpDataInDatabase := model.UsersOtpData{}

	db.WithContext(ctx).Where(
		model.UsersOtpData{
			UserId: ctx.Value(constant.UserIdLogKey).(string),
		},
	).Find(&userOtpDataInDatabase)

	// if so, check the input otp before deleting the file or folder
	if userOtpDataInDatabase.BaseEntity.Id != 0 {
		if inputOtp == "" {
			return constant.InputOtpEmptyError, errors.New("otp is empty")
		}
		userCurrentOtp, otpGeneratorError := GenerateTotp(ctx, userOtpDataInDatabase.UserOtpSecretData)
		if otpGeneratorError != nil {
			return constant.OtpError, otpGeneratorError
		}
		log.WithLevel(constant.Info, ctx, "Current OTP is %s", userCurrentOtp)
		if userCurrentOtp != inputOtp {
			return constant.WrongOtpError, errors.New("check your otp")
		}
	}
	return constant.Success, nil
}
