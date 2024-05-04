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

type Tabler interface {
	TableName() string
}

func (UserAccountActivationLog) TableName() string {
	return "user_account_activation_log"
}
