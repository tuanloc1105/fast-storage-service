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

func KeycloakRevokeToken(ctx context.Context, refreshToken string) (payload.RevokeTokenError, error) {
	revokeTokenCurlCommand := fmt.Sprintf(
		"curl -k --location '%s' --header 'Content-Type: application/x-www-form-urlencoded' --data-urlencode 'client_id=%s' --data-urlencode 'client_secret=%s' --data-urlencode 'token=%s'",
		config.KeycloakApiUrl+config.KeycloakRevokeTokenPath,
		config.KeycloakClientId,
		config.KeycloakClientSecret,
		refreshToken,
	)
	result := payload.RevokeTokenError{}
	shellStdout, _, shellError := utils.Shellout(ctx, revokeTokenCurlCommand)
	if shellError != nil {
		log.WithLevel(constant.Info, ctx, "an error has been occurred: %s", shellError.Error())
		return result, shellError
	}

	utils.JsonToStruct(shellStdout, &result)
	return result, nil
}
