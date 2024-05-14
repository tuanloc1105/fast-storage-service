package keycloak

import (
	"context"
	"errors"
	"fast-storage-go-service/config"
	"fast-storage-go-service/payload"
)

func KeycloakDeactivateAccount(ctx context.Context, username string) error {
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
				return errors.New("can not deactivate user")
			}
			user := userSearchingResult[0]
			updateUserStatus := payload.RealmsUsersCreationInfo{
				Username:      user.Username,
				Email:         user.Email,
				FirstName:     user.FirstName,
				LastName:      user.LastName,
				EmailVerified: true,
				Enabled:       false,
			}
			return KeycloakUpdateUser(ctx, adminProtocolOpenidConnectToken.AccessToken, user.ID, updateUserStatus)

		}
	}
}
