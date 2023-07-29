package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID                  uuid.UUID `gorm:"primary_key"`
	KTP                 string    `gorm:"not null"`
	Phonenumber         string    `gorm:"not null"`
	Gender              int       `gorm:"not null"`
	BirthDate           string    `gorm:"not null"`
	Name                string    `gorm:"type:varchar(255);not null" `
	Email               string    `gorm:"uniqueIndex;not null"`
	PinNumber           string    `gorm:"null"`
	LimitCredit         uint64    `gorm:"default:1000000"`
	TotalCredit         uint64    `gorm:"default:0"`
	Password            string    `gorm:"not null"`
	AccountName         string    `gorm:"null"`
	AccountNumber       string    `gorm:"null"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	SentTransaction     []Transaction `gorm:"foreignKey:SenderID"`
	ReceivedTransaction []Transaction `gorm:"foreignKey:ReceiverID"`
}

type UserModel struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name" binding:"required"`
	KTP           string    `json:"ktp"`
	Phonenumber   string    `json:"phonenumber"`
	BirthDate     string    `json:"birt_date"`
	Gender        int       `json:"gender"`
	Email         string    `json:"email" binding:"required"`
	Password      string    `json:"password" binding:"required"`
	PinNumber     string    `json:"pin_number"`
	LimitCredit   uint64    `json:"limit_credit"`
	TotalCredit   uint64    `json:"total_creidt"`
	AccountName   string    `json:"acount_name"`
	AccountNumber string    `json:"acount_number"`
	Balance       int       `json:"balance"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
