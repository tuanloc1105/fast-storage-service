package payload

type CutOrCopyRequestBodyValue struct {
	SourceFolder      string `json:"sourceFolder"`
	DestinationFolder string `json:"destinationFolder"`
	FileName          string `json:"fileName"`
	IsCopy            bool   `json:"isCopy"`
}

type CutOrCopyRequestBody struct {
	Request CutOrCopyRequestBodyValue `json:"request" binding:"required"`
}

type ReadImageFileRequestBodyValue struct {
	FolderLocation string `json:"folderLocation"`
	ImageFileName  string `json:"imageFileName"`
}

type ReadImageFileRequestBody struct {
	Request ReadImageFileRequestBodyValue `json:"request" binding:"required"`
}

type SearchFileRequestBodyValue struct {
	SearchingContent string `json:"searchingContent"`
}

type SearchFileRequestBody struct {
	Request SearchFileRequestBodyValue `json:"request" binding:"required"`
}

type EncryptEveryFolderRequestBodyValue struct {
	Encrypt bool `json:"encrypt" binding:"required"`
}

type EncryptEveryFolderRequestBody struct {
	Request EncryptEveryFolderRequestBodyValue `json:"request" binding:"required"`
}

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
	Credential      string `json:"credential"`
}

type GetAllElementInSpecificDirectoryBody struct {
	Request GetAllElementInSpecificDirectoryBodyValue `json:"request" binding:"required"`
}

type RemoveFileBodyValue struct {
	LocationToRemove string   `json:"locationToRemove" binding:"required"`
	FileNameToRemove []string `json:"fileNameToRemove"`
	OtpCredential    string   `json:"otpCredential"`
	Credential       string   `json:"credential"`
}

type RemoveFileBody struct {
	Request RemoveFileBodyValue `json:"request" binding:"required"`
}

type CreateFolderBodyValue struct {
	FolderToCreate string `json:"folderToCreate" binding:"required"`
}

type CreateFolderBody struct {
	Request CreateFolderBodyValue `json:"request" binding:"required"`
}

type RenameFolderBodyValue struct {
	OldFolderLocationName string `json:"oldFolderLocationName" binding:"required"`
	NewFolderLocationName string `json:"newFolderLocationName" binding:"required"`
}

type RenameFolderBody struct {
	Request RenameFolderBodyValue `json:"request" binding:"required"`
}

type CreateFileBodyValue struct {
	FolderToCreate   string `json:"folderToCreate" binding:"required"`
	FileNameToCreate string `json:"fileNameToCreate" binding:"required"`
	FileExtension    string `json:"fileExtension" binding:"required"`
}

type CreateFileBody struct {
	Request CreateFileBodyValue `json:"request" binding:"required"`
}

type SetPasswordForFolderBodyValue struct {
	Folder         string `json:"folder" binding:"required"`
	CredentialType string `json:"credentialType" binding:"required"`
	Credential     string `json:"credential"`
}

type SetPasswordForFolderBody struct {
	Request SetPasswordForFolderBodyValue `json:"request" binding:"required"`
}

type CheckSecureFolderStatusBodyValue struct {
	Folder string `json:"folder" binding:"required"`
}

type CheckSecureFolderStatusBody struct {
	Request CheckSecureFolderStatusBodyValue `json:"request" binding:"required"`
}

type ShareFileBodyValue struct {
	Folder                             string   `json:"folder" binding:"required"`
	File                               string   `json:"file" binding:"required"`
	UserEmailToShare                   []string `json:"userEmailToShare" binding:"required"`
	TheTimeIntervalInMinutesToBeShared int      `json:"theTimeIntervalInMinutesToBeShared" binding:"required"`
}

type ShareFileBody struct {
	Request ShareFileBodyValue `json:"request" binding:"required"`
}

type DownloadMultipleFileBodyValue struct {
	Folder     string `json:"folder" binding:"required"`
	File       string `json:"file" binding:"required"`
	Credential string `json:"credential"`
}

type DownloadMultipleFileBody struct {
	Request DownloadMultipleFileBodyValue `json:"request" binding:"required"`
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
	Size             string `json:"size"`
	Name             string `json:"name"`
	Extension        string `json:"extension"`
	LastModifiedDate string `json:"lastModifiedDate"`
	Type             string `json:"type"`
	Editable         bool   `json:"editable"`
	BirthDate        string `json:"birthDate"`
}

type ImageViewResponse struct {
	Data string `json:"data"`
}
