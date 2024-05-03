package keycloak

import (
	"context"
	"fast-storage-go-service/config"
	"fast-storage-go-service/constant"
	"fast-storage-go-service/log"
	"fast-storage-go-service/payload"
	"fast-storage-go-service/utils"
	"fmt"
	"net/url"
)

func KeycloakSearchUser(ctx context.Context, adminAccessToken string, searchContent string) ([]payload.RealmsUsersInfo, error) {
	searchUserCurlCommand := fmt.Sprintf(
		"curl -k --location '%s' --header 'Authorization: Bearer %s'",
		fmt.Sprintf(config.KeycloakApiUrl+config.KeycloakSearchUserPath, url.QueryEscape(searchContent)),
		adminAccessToken,
	)

	result := []payload.RealmsUsersInfo{}
	shellStdout, _, shellError := utils.Shellout(ctx, searchUserCurlCommand)
	if shellError != nil {
		log.WithLevel(constant.Info, ctx, "an error has been occurred: %s", shellError.Error())
		return result, shellError
	}

	utils.JsonToStruct(shellStdout, &result)
	return result, nil
}
