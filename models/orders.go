package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	ID        uint64  `gorm:"primary_key"`
	UserID    string  `gorm:"type:varchar(255);not null" `
	Amount    float32 `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type OrderModel struct {
	ID        uint64    `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
