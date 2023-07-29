package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Bill struct {
	gorm.Model
	ID       uint64    `gorm:"primary_key"`
	UserID   uuid.UUID `gorm:"not null"`
	Sender   User      `gorm:"foreignKey:SenderID"`
	Receiver User      `gorm:"foreignKey:ReceiverID"`
}
