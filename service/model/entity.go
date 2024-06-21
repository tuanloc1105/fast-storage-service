package model

import (
	"time"
)

type BaseEntity struct {
	Id        int64     `json:"id" gorm:"column:id;primaryKey;not null"`
	UUID      string    `json:"uuid" gorm:"column:uuid;not null"`
	Active    *bool     `json:"active" gorm:"column:active;not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;not null"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;not null"`
	CreatedBy string    `json:"createdBy" gorm:"column:created_by;not null"`
	UpdatedBy string    `json:"updatedBy" gorm:"column:updated_by;not null"`
}

type UserAccountActivationLog struct {
	BaseEntity BaseEntity `gorm:"embedded" json:"baseInfo"`
	UserId     string     `json:"userId" gorm:"column:user_id;not null;unique"`
	Username   string     `json:"username" gorm:"column:username;not null"`
}

type UsersOtpData struct {
	BaseEntity                   BaseEntity `gorm:"embedded" json:"baseInfo"`
	UserId                       string     `json:"userId" gorm:"column:user_id;not null;unique"`
	UserOtpSecretData            string     `json:"userOtpSecretData" gorm:"column:user_otp_secret_data;not null"`
	UserOtpQrCodeImageBase64Data string     `json:"userOtpQrCodeImageBase64Data" gorm:"column:user_otp_qr_code_image_base64_data;not null"`
}

type UserAuthenticationLog struct {
	BaseEntity                     BaseEntity `gorm:"embedded" json:"baseInfo"`
	Username                       string     `json:"username" gorm:"column:Username;not null"`
	AuthenticatedAt                time.Time  `json:"authenticatedAt" gorm:"column:authenticated_at;not null"`
	AuthenticatedStatus            string     `json:"authenticatedStatus" gorm:"column:authenticated_status;not null"`
	AuthenticatedStatusDescription string     `json:"authenticatedStatusDescription" gorm:"column:authenticated_status_description"`
}

type UserStorageLimitationData struct {
	BaseEntity         BaseEntity `gorm:"embedded" json:"baseInfo"`
	Username           string     `json:"username" gorm:"column:username;not null;unique"`
	MaximunStorageSize int        `json:"maximunStorageSize" gorm:"column:maximun_storage_size;not null"`
	StorageSizeUnit    string     `json:"storageSizeUnit" gorm:"column:storage_size_unit;not null"`
}

type UserFolderCredential struct {
	BaseEntity               BaseEntity `gorm:"embedded" json:"baseInfo"`
	Username                 string     `json:"username" gorm:"column:username;not null"`
	Directory                string     `json:"directory" gorm:"column:directory;not null"`
	Credential               string     `json:"credential" gorm:"column:credential;not null"`                                // bcrypted password, empty string if credential type is OTP
	CredentialType           string     `json:"credentialType" gorm:"column:credential_type;not null"`                       // OTP or PASSWORD
	LastFolderActivitiesTime time.Time  `json:"lastFolderActivitiesTime" gorm:"column:last_folder_activities_time;not null"` // folder idle time is 5 minutes last, if exceed then user have to be enter password again to view folder content
}

type UserFileAndFolderSharing struct {
	BaseEntity        BaseEntity `gorm:"embedded" json:"baseInfo"`
	Username          string     `json:"username" gorm:"column:username;not null"`
	ListOfUsersShared string     `json:"listOfUsersShared" gorm:"column:list_of_users_shared"`
	Directory         string     `json:"directory" gorm:"column:directory;not null"`
	FileName          string     `json:"fileName" gorm:"column:file_name"`
	ExpiredTime       string     `json:"expiredTime" gorm:"column:expired_time"`
}

// type Test struct {
// 	BaseEntity                     `gorm:"embedded" json:"baseInfo"`
// 	Username                       string    `json:"username" gorm:"column:Username;not null"`
// 	AuthenticatedAt                time.Time `json:"authenticatedAt" gorm:"column:authenticated_at;not null"`
// 	AuthenticatedStatus            string    `json:"authenticatedStatus" gorm:"column:authenticated_status;not null"`
// 	AuthenticatedStatusDescription string    `json:"authenticatedStatusDescription" gorm:"column:authenticated_status_description"`
// }

type Tabler interface {
	TableName() string
}

func (UserAccountActivationLog) TableName() string {
	return "user_account_activation_log"
}

func (UsersOtpData) TableName() string {
	return "users_otp_data"
}

func (UserAuthenticationLog) TableName() string {
	return "user_authentication_log"
}

func (UserStorageLimitationData) TableName() string {
	return "user_storage_limitation_data"
}

func (UserFolderCredential) TableName() string {
	return "user_folder_credential"
}

func (UserFileAndFolderSharing) TableName() string {
	return "user_file_and_folder_sharing"
}
