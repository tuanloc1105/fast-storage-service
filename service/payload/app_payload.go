package payload

type RegisterRequestBodyValue struct {
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
	Email           string `json:"email" binding:"required"`
	FirstName       string `json:"firstName" binding:"required"`
	LastName        string `json:"lastName" binding:"required"`
}

type RegisterRequestBody struct {
	Request RegisterRequestBodyValue `json:"request" binding:"required"`
}

type LoginRequestBodyValue struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequestBody struct {
	Request LoginRequestBodyValue `json:"request" binding:"required"`
}

type GetNewTokenBodyValue struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type GetNewTokenBody struct {
	Request GetNewTokenBodyValue `json:"request" binding:"required"`
}

type LogoutBodyValue struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type LogoutBody struct {
	Request LogoutBodyValue `json:"request" binding:"required"`
}

type GetAllElementInSpecificDirectoryBodyValue struct {
	CurrentLocation string `json:"currentLocation"`
}

type GetAllElementInSpecificDirectoryBody struct {
	Request GetAllElementInSpecificDirectoryBodyValue `json:"request" binding:"required"`
}

type DownloadFileBodyValue struct {
	LocationToDownload string `json:"locationToDownload" binding:"required"`
	FileNameToDownload string `json:"fileNameToDownload" binding:"required"`
}

type DownloadFileBody struct {
	Request DownloadFileBodyValue `json:"request" binding:"required"`
}

type RemoveFileBodyValue struct {
	LocationToRemove string `json:"locationToRemove" binding:"required"`
	FileNameToRemove string `json:"fileNameToRemove"`
	OtpCredential    string `json:"otpCredential"`
}

type RemoveFileBody struct {
	Request RemoveFileBodyValue `json:"request" binding:"required"`
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

type KeycloakCommonErrorResponseResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"errorDescription"`
	ErrorMessage     string `json:"errorMessage"`
}

type SystemStorageStatus struct {
	Size            string `json:"size"`
	Used            string `json:"used"`
	Avail           string `json:"avail"`
	UseInPercentage string `json:"useInPercentage"`
}

type UserStorageStatus struct {
	MaximunSize float64 `json:"maximunSize"`
	Used        float64 `json:"used"`
}

type FileInformation struct {
	Size string `json:"size"`
	Name string `json:"name"`
	Type string `json:"type"`
}
