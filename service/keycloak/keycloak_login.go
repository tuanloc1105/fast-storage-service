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

func KeycloakLogin(ctx context.Context, username string, password string) (payload.ProtocolOpenidConnectToken, error) {

	loginCurlCommand := fmt.Sprintf(
		"curl -k --connect-timeout 30 --max-time 40 --location '%s' --header 'Content-Type: application/x-www-form-urlencoded' --data-urlencode 'client_id=%s' --data-urlencode 'client_secret=%s' --data-urlencode 'username=%s' --data-urlencode 'password=%s' --data-urlencode 'grant_type=password' --data-urlencode 'scope=openid'",
		config.KeycloakApiUrl+config.KeycloakLoginPath,
		config.KeycloakClientId,
		config.KeycloakClientSecret,
		username,
		password,
	)
	result := payload.ProtocolOpenidConnectToken{}
	shellStdout, _, shellError := utils.Shellout(ctx, loginCurlCommand)
	if shellError != nil {
		log.WithLevel(constant.Info, ctx, "an error has been occurred: %s", shellError.Error())
		return result, shellError
	}

	utils.JsonToStruct(shellStdout, &result)
	return result, nil
}
