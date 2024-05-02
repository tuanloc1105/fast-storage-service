package config

import "os"

var (
	KeycloakClientId     = ""
	KeycloakClientSecret = ""
	KeycloakApiUrl       = ""
)

func CheckKeycloakInfo() bool {
	keycloakClientId, isKeycloakClientIdSet := os.LookupEnv("KEYCLOAK_CLIENT_ID")
	keycloakClientSecret, isKeycloakClientSecretSet := os.LookupEnv("KEYCLOAK_CLIENT_SECRET")
	keycloakApiUrl, isKeycloakApiUrlSet := os.LookupEnv("KEYCLOAK_API_URL")

	KeycloakClientId = keycloakClientId
	KeycloakClientSecret = keycloakClientSecret
	KeycloakApiUrl = keycloakApiUrl

	return isKeycloakClientIdSet && isKeycloakClientSecretSet && isKeycloakApiUrlSet

}
