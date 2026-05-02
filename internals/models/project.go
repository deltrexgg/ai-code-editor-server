package models

import (
	"time"

	"github.com/google/uuid"
)

type Projects struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuidv7()" json:"id"`

	UserID	uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`

	ProjectName	string `gorm:"type:text;not null;index;" json:"project_name"`
	Description	string `gorm:"type:text" json:"description"`
	TechStack	string `gorm:"type:text" json:"tech_stack"`

	IsPublished	bool `gorm:"default:false" json:"is_published"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}