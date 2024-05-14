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
