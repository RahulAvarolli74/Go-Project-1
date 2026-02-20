package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Recipe struct {
	ID            string    `gorm:"type:text;primaryKey" json:"id"`
	Title         string    `gorm:"type:text;not null" json:"title" binding:"required"`
	Description   string    `gorm:"type:text" json:"description"`
	ImageURL      string    `gorm:"type:text" json:"image_url"`
	Ingredients   string    `gorm:"type:text" json:"ingredients" binding:"required"`
	PrepTime      int       `gorm:"default:0" json:"prep_time"`
	CookTime      int       `gorm:"default:0" json:"cook_time"`
	Servings      int       `gorm:"default:1" json:"servings"`
	AverageRating float64   `gorm:"default:0" json:"average_rating"`
	UserID        string    `gorm:"type:text;index" json:"user_id"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Ratings       []Rating  `gorm:"foreignKey:RecipeID" json:"ratings,omitempty"`
}

func (r *Recipe) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}
