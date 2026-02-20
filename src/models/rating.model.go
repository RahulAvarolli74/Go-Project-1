package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Rating struct {
	ID        string    `gorm:"type:text;primaryKey" json:"id"`
	RecipeID  string    `gorm:"type:text;index;not null" json:"recipe_id"`
	UserName  string    `gorm:"type:text;not null" json:"user_name" binding:"required"`
	Score     int       `gorm:"not null" json:"score" binding:"required,min=1,max=5"`
	Comment   string    `gorm:"type:text" json:"comment"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (r *Rating) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}
