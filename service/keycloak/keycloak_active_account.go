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
)

func KeycloakActiveAccount(ctx context.Context, userId, username string) error {

	// login with admin to register user
	if adminProtocolOpenidConnectToken, protocolOpenidConnectTokenError := KeycloakLogin(ctx, config.KeycloakAdminUsername, config.KeycloakAdminPassword); protocolOpenidConnectTokenError != nil {
		return protocolOpenidConnectTokenError
	} else {

		if adminProtocolOpenidConnectToken.Error != "" {
			return errors.New(adminProtocolOpenidConnectToken.ErrorDescription)
		}

		if userSearchingResult, userSearchingResultError := KeycloakSearchUser(ctx, adminProtocolOpenidConnectToken.AccessToken, username); userSearchingResultError != nil {
			return userSearchingResultError
		} else {
			if len(userSearchingResult) != 1 {
				return errors.New("can not create user")
			}
			user := userSearchingResult[0]
			realmsUsersCreationInfo := payload.RealmsUsersCreationInfo{
				Username:      user.Username,
				Email:         user.Email,
				FirstName:     user.FirstName,
				LastName:      user.LastName,
				EmailVerified: true,
				Enabled:       true,
			}
			updateUserJsonString := utils.StructToJson(realmsUsersCreationInfo)
			userRegisterCurlCommand := fmt.Sprintf(
				"curl -k --location --request PUT '%s' --header 'Content-Type: application/json' --header 'Authorization: Bearer %s' --data-raw '%s'",
				fmt.Sprintf(config.KeycloakApiUrl+config.KeycloakUpdateUserPath, userId),
				adminProtocolOpenidConnectToken.AccessToken,
				updateUserJsonString,
			)
			shellStdout, _, shellError := utils.Shellout(ctx, userRegisterCurlCommand)
			if shellError != nil {
				log.WithLevel(constant.Info, ctx, "an error has been occurred: %s", shellError.Error())
				return shellError
			}

			keycloakCommonError := payload.KeycloakCommonErrorResponse{}
			utils.JsonToStruct(shellStdout, &keycloakCommonError)

			if keycloakCommonError.Error != "" {
				return errors.New(keycloakCommonError.ErrorDescription)
			}

			if keycloakCommonError.ErrorMessage != "" {
				return errors.New(keycloakCommonError.ErrorMessage)
			}
		}
	}
	return nil
}
