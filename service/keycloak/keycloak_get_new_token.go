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

func KeycloakGetNewToken(ctx context.Context, refreshToken string) (payload.ProtocolOpenidConnectToken, error) {
	getNewTokenCurlCommand := fmt.Sprintf(
		"curl -k --connect-timeout 30 --max-time 40 --location '%s' --header 'Content-Type: application/x-www-form-urlencoded' --data-urlencode 'client_id=%s' --data-urlencode 'client_secret=%s' --data-urlencode 'grant_type=refresh_token' --data-urlencode 'refresh_token=%s'",
		config.KeycloakApiUrl+config.KeycloakGetNewTokenPath,
		config.KeycloakClientId,
		config.KeycloakClientSecret,
		refreshToken,
	)
	result := payload.ProtocolOpenidConnectToken{}
	shellStdout, _, shellError := utils.Shellout(ctx, getNewTokenCurlCommand)
	if shellError != nil {
		log.WithLevel(constant.Info, ctx, "an error has been occurred: %s", shellError.Error())
		return result, shellError
	}

	utils.JsonToStruct(shellStdout, &result)
	return result, nil
}
