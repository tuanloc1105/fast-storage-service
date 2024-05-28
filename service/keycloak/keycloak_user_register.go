package keycloak

import (
	"context"
	"errors"
	"fast-storage-go-service/config"
	"fast-storage-go-service/constant"
	"fast-storage-go-service/log"
	"fast-storage-go-service/payload"
	"fast-storage-go-service/utils"
	"fmt"
	"os"
	"strings"
)

func KeycloakUserRegister(ctx context.Context, input payload.RegisterRequestBodyValue) (string, error) {

	if strings.Compare(input.Password, input.ConfirmPassword) != 0 {
		return "", errors.New("password and confirm password is not same")
	}

	// login with admin to register user
	if adminProtocolOpenidConnectToken, protocolOpenidConnectTokenError := KeycloakLogin(ctx, config.KeycloakAdminUsername, config.KeycloakAdminPassword); protocolOpenidConnectTokenError != nil {
		return "", protocolOpenidConnectTokenError
	} else {

		if adminProtocolOpenidConnectToken.Error != "" {
			return "", errors.New(adminProtocolOpenidConnectToken.ErrorDescription)
		}

		// register user
		registerUserRequest := payload.RealmsUsersCreationInfo{
			Username:      input.Username,
			Email:         input.Email,
			FirstName:     input.FirstName,
			LastName:      input.LastName,
			EmailVerified: false,
			Enabled:       false, // active thourgh another api
		}
		registerUserRequestJsonString := utils.StructToJson(registerUserRequest)
		userRegisterCurlCommand := fmt.Sprintf(
			"curl -k --connect-timeout 30 --max-time 40 --location '%s' --header 'Content-Type: application/json' --header 'Authorization: Bearer %s' --data-raw '%s'",
			config.KeycloakApiUrl+config.KeycloakUserRegisterPath,
			adminProtocolOpenidConnectToken.AccessToken,
			registerUserRequestJsonString,
		)
		shellStdout, _, shellError := utils.Shellout(ctx, userRegisterCurlCommand)
		if shellError != nil {
			log.WithLevel(constant.Info, ctx, "an error has been occurred: %s", shellError.Error())
			return "", shellError
		}

		keycloakCommonError := payload.KeycloakCommonErrorResponse{}
		utils.JsonToStruct(shellStdout, &keycloakCommonError)

		if keycloakCommonError.Error != "" {
			return "", errors.New(keycloakCommonError.ErrorDescription)
		}

		if keycloakCommonError.ErrorMessage != "" {
			return "", errors.New(keycloakCommonError.ErrorMessage)
		}

		var activeLink string

		// get user id after created
		if userSearchingResult, userSearchingResultError := KeycloakSearchUser(ctx, adminProtocolOpenidConnectToken.AccessToken, registerUserRequest.Email); userSearchingResultError != nil {
			return "", userSearchingResultError
		} else {
			if len(userSearchingResult) != 1 {
				return "", errors.New("can not create user")
			}
			// set password for user
			resetPasswordInput := payload.ResetPasswordKeycloakInput{
				Temporary: false,
				Type:      "password",
				Value:     input.Password,
			}

			resetPasswordCurlCommand := fmt.Sprintf(
				"curl -k --connect-timeout 30 --max-time 40 --location --request PUT '%s' --header 'Content-Type: application/json' --header 'Authorization: Bearer %s' --data '%s'",
				fmt.Sprintf(config.KeycloakApiUrl+config.KeycloakResetPasswordPath, userSearchingResult[0].ID),
				adminProtocolOpenidConnectToken.AccessToken,
				utils.StructToJson(resetPasswordInput),
			)
			resetPasswordShellStdout, _, resetPasswordShellError := utils.Shellout(ctx, resetPasswordCurlCommand)
			if resetPasswordShellError != nil {
				log.WithLevel(constant.Info, ctx, "an error has been occurred: %s", resetPasswordShellError.Error())
				return "", resetPasswordShellError
			}
			utils.JsonToStruct(resetPasswordShellStdout, &keycloakCommonError)

			if keycloakCommonError.Error != "" {
				return "", errors.New(keycloakCommonError.ErrorDescription)
			}

			if keycloakCommonError.ErrorMessage != "" {
				return "", errors.New(keycloakCommonError.ErrorMessage)
			}

			activeHost, activeHostSet := os.LookupEnv("ACCOUNT_ACTIVE_HOST")

			if !activeHostSet {
				activeHost = "http://localhost:8080"
			}

			activeLink = fmt.Sprintf("%s/fast_storage/api/v1/auth/active_account?userId=%s&username=%s", activeHost, userSearchingResult[0].ID, userSearchingResult[0].Username)
			log.WithLevel(constant.Info, ctx, "active link for user: %s", activeLink)
		}

		return activeLink, nil
	}
}
