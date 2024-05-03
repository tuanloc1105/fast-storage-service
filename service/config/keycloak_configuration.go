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
	KeycloakLoginPath       = "/realms/master/protocol/openid-connect/token"
	KeycloakGetUserInfoPath = "/realms/master/protocol/openid-connect/token/introspect"
	KeycloakGetNewTokenPath = "/realms/master/protocol/openid-connect/token"
	KeycloakRevokeTokenPath = "/realms/master/protocol/openid-connect/revoke"
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
