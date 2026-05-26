package user

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint   `gorm:"primaryKey" json:"id" example:"1"`
	Username     string `gorm:"uniqueIndex;not null;type:varchar(50)" json:"username" example:"adam_99"`
	DisplayName  string `gorm:"not null; type:varchar(100)" json:"display_name" example:"BlazeDev"`
	PasswordHash string `gorm:"not null" json:"-"`
	CreatedAt    time.Time `json:"created_at" example:"2023-10-27T10:00:00Z"`
	UpdatedAt    time.Time `json:"updated_at" example:"2023-10-27T10:00:00Z"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}