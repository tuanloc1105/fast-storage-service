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
