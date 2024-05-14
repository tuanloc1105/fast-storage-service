package keycloak

import (
	"context"
	"fast-storage-go-service/config"
	"fast-storage-go-service/constant"
	"fast-storage-go-service/log"
	"fast-storage-go-service/payload"
	"fast-storage-go-service/utils"
	"fmt"
)

func KeycloakGetUserInfo(ctx context.Context, token string) (payload.OpenidConnectTokenIntrospect, error) {
	getUserInfoCurlCommand := fmt.Sprintf(
		"curl -k --connect-timeout 30 --max-time 40 --location '%s' --header 'Content-Type: application/x-www-form-urlencoded' --data-urlencode 'client_id=%s' --data-urlencode 'client_secret=%s' --data-urlencode 'token=%s'",
		config.KeycloakApiUrl+config.KeycloakGetUserInfoPath,
		config.KeycloakClientId,
		config.KeycloakClientSecret,
		token,
	)
	result := payload.OpenidConnectTokenIntrospect{}
	shellStdout, _, shellError := utils.Shellout(ctx, getUserInfoCurlCommand)
	if shellError != nil {
		log.WithLevel(constant.Info, ctx, "an error has been occurred: %s", shellError.Error())
		return result, shellError
	}

	utils.JsonToStruct(shellStdout, &result)
	return result, nil
}
