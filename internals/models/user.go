package models

import (
	"time"

	"github.com/google/uuid"
)

type Users struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuidv7()" json:"id"`

	Name		string `gorm:"type:text;not null;index" json:"name"`
	Email		string `gorm:"type:text;not null;uniqueIndex;index" json:"email"`
	Password	string `gorm:"type:text;not null" json:"password,omitempty"`

	Designation	string `gorm:"type:text" json:"designation"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}