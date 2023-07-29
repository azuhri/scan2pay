package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	ID              uint64    `gorm:"primary_key"`
	TransactionCode string    `gorm:"uniqueIndex; not null"`
	SenderID        uuid.UUID `gorm:"not null"`
	ReceiverID      uuid.UUID `gorm:"not null"`
	Amount          int       `gorm:"not null"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Sender          User `gorm:"foreignKey:SenderID"`
	Receiver        User `gorm:"foreignKey:ReceiverID"`
}
