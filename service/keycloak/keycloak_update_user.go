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

func KeycloakUpdateUser(ctx context.Context, administratorAccessToken, userId string, requestBody payload.RealmsUsersCreationInfo) error {
	requestBodyString := utils.StructToJson(requestBody)
	updateUserCurlCommand := fmt.Sprintf(
		"curl -k --connect-timeout 30 --max-time 40 --location --request PUT '%s' --header 'Content-Type: application/json' --header 'Authorization: Bearer %s' --data-raw '%s'",
		fmt.Sprintf(config.KeycloakApiUrl+config.KeycloakUpdateUserPath, userId),
		administratorAccessToken,
		requestBodyString,
	)
	shellStdout, _, shellError := utils.Shellout(ctx, updateUserCurlCommand)
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

	return nil
}
