package config

import "os"

var (
	KeycloakClientId      = ""
	KeycloakClientSecret  = ""
	KeycloakApiUrl        = ""
	KeycloakAdminUsername = ""
	KeycloakAdminPassword = ""
)

const (
	KeycloakRealm             = "master"
	KeycloakLoginPath         = "/realms/" + KeycloakRealm + "/protocol/openid-connect/token"
	KeycloakGetUserInfoPath   = "/realms/" + KeycloakRealm + "/protocol/openid-connect/token/introspect"
	KeycloakGetNewTokenPath   = "/realms/" + KeycloakRealm + "/protocol/openid-connect/token"
	KeycloakRevokeTokenPath   = "/realms/" + KeycloakRealm + "/protocol/openid-connect/revoke"
	KeycloakUserRegisterPath  = "/admin/realms/" + KeycloakRealm + "/users"
	KeycloakSearchUserPath    = "/admin/realms/" + KeycloakRealm + "/users?briefRepresentation=true&search=%s"
	KeycloakResetPasswordPath = "/admin/realms/" + KeycloakRealm + "/users/%s/reset-password"
)

func CheckKeycloakInfo() bool {
	keycloakClientId, isKeycloakClientIdSet := os.LookupEnv("KEYCLOAK_CLIENT_ID")
	keycloakClientSecret, isKeycloakClientSecretSet := os.LookupEnv("KEYCLOAK_CLIENT_SECRET")
	keycloakApiUrl, isKeycloakApiUrlSet := os.LookupEnv("KEYCLOAK_API_URL")
	keycloakAdminUsername, keycloakAdminUsernameSet := os.LookupEnv("KEYCLOAK_ADMIN_USERNAME")
	keycloakAdminPassword, keycloakAdminPasswordSet := os.LookupEnv("KEYCLOAK_ADMIN_PASSWORD")

	if !keycloakAdminUsernameSet {
		keycloakAdminUsername = "admin"
	}
	if !keycloakAdminPasswordSet {
		keycloakAdminPassword = "admin"
	}

	KeycloakClientId = keycloakClientId
	KeycloakClientSecret = keycloakClientSecret
	KeycloakApiUrl = keycloakApiUrl
	KeycloakAdminUsername = keycloakAdminUsername
	KeycloakAdminPassword = keycloakAdminPassword

	return isKeycloakClientIdSet && isKeycloakClientSecretSet && isKeycloakApiUrlSet

}
