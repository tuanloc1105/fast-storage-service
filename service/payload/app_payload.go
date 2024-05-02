package payload

type LoginRequestBodyValue struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequestBody struct {
	Request LoginRequestBodyValue `json:"request"`
}

type ProtocolOpenidConnectTokenResponse struct {
	AccessToken      string `json:"accessToken"`
	ExpiresIn        int64  `json:"expiresIn"`
	RefreshExpiresIn int64  `json:"refreshExpiresIn"`
	RefreshToken     string `json:"refreshToken"`
	TokenType        string `json:"tokenType"`
	IDToken          string `json:"idToken"`
	NotBeforePolicy  int64  `json:"notBeforePolicy"`
	SessionState     string `json:"sessionState"`
	Scope            string `json:"scope"`
	Error            string `json:"error"`
	ErrorDescription string `json:"errorDescription"`
}

type OpenidConnectTokenIntrospectResponse struct {
	Exp               int64                  `json:"exp"`
	Iat               int64                  `json:"iat"`
	Jti               string                 `json:"jti"`
	Iss               string                 `json:"iss"`
	Aud               []string               `json:"aud"`
	Sub               string                 `json:"sub"`
	Typ               string                 `json:"typ"`
	Azp               string                 `json:"azp"`
	SessionState      string                 `json:"sessionState"`
	ACR               string                 `json:"acr"`
	AllowedOrigins    []string               `json:"allowedOrigins"`
	RealmAccess       RealmAccessResponse    `json:"realmAccess"`
	ResourceAccess    ResourceAccessResponse `json:"resourceAccess"`
	Scope             string                 `json:"scope"`
	Sid               string                 `json:"sid"`
	EmailVerified     bool                   `json:"emailVerified"`
	PreferredUsername string                 `json:"preferredUsername"`
	ClientID          string                 `json:"clientId"`
	Username          string                 `json:"username"`
	TokenType         string                 `json:"tokenType"`
	Active            bool                   `json:"active"`
}

type RealmAccessResponse struct {
	Roles []string `json:"roles"`
}

type ResourceAccessResponse struct {
	MasterRealm RealmAccessResponse `json:"masterRealm"`
	Account     RealmAccessResponse `json:"account"`
}
