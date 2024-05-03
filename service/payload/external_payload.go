package payload

type ProtocolOpenidConnectToken struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int64  `json:"expires_in"`
	RefreshExpiresIn int64  `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	IDToken          string `json:"id_token"`
	NotBeforePolicy  int64  `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type OpenidConnectTokenIntrospect struct {
	Exp               int64          `json:"exp"`
	Iat               int64          `json:"iat"`
	Jti               string         `json:"jti"`
	Iss               string         `json:"iss"`
	Aud               []string       `json:"aud"`
	Sub               string         `json:"sub"`
	Typ               string         `json:"typ"`
	Azp               string         `json:"azp"`
	SessionState      string         `json:"session_state"`
	ACR               string         `json:"acr"`
	AllowedOrigins    []string       `json:"allowed-origins"`
	RealmAccess       RealmAccess    `json:"realm_access"`
	ResourceAccess    ResourceAccess `json:"resource_access"`
	Scope             string         `json:"scope"`
	Sid               string         `json:"sid"`
	EmailVerified     bool           `json:"email_verified"`
	PreferredUsername string         `json:"preferred_username"`
	ClientID          string         `json:"client_id"`
	Username          string         `json:"username"`
	TokenType         string         `json:"token_type"`
	Active            bool           `json:"active"`
}

type RealmAccess struct {
	Roles []string `json:"roles"`
}

type ResourceAccess struct {
	MasterRealm RealmAccess `json:"master-realm"`
	Account     RealmAccess `json:"account"`
}

type RevokeTokenError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
