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

func KeycloakGetUserInfo(ctx context.Context, token string) (payload.RealmsUsersInfoElement, error) {
	getUserInfoCurlCommand := fmt.Sprintf(
		"curl -k --location '%s/admin/realms/master/users/' --header 'Authorization: Bearer %s'",
		config.KeycloakApiUrl,
		token,
	)
	result := []payload.RealmsUsersInfoElement{}
	shellStdout, _, shellError := utils.Shellout(ctx, getUserInfoCurlCommand)
	if shellError != nil {
		log.WithLevel(constant.Info, ctx, "an error has been occurred: %s", shellError.Error())
		return payload.RealmsUsersInfoElement{}, shellError
	}

	utils.JsonToStruct(shellStdout, &result)
	if len(result) < 1 {
		return payload.RealmsUsersInfoElement{}, nil
	}
	return result[0], nil
}
