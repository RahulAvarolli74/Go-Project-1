package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        string    `gorm:"type:text;primaryKey" json:"id"`
	Username  string    `gorm:"type:text;uniqueIndex;not null" json:"username" binding:"required"`
	Email     string    `gorm:"type:text;uniqueIndex;not null" json:"email" binding:"required,email"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Recipes   []Recipe  `gorm:"foreignKey:UserID" json:"recipes,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}
