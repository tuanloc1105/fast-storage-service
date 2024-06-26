package keycloak

import (
	"context"
	"errors"
	"fast-storage-go-service/config"
	"fast-storage-go-service/payload"
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
			return KeycloakUpdateUser(ctx, adminProtocolOpenidConnectToken.AccessToken, userId, realmsUsersCreationInfo)
		}
	}
}
