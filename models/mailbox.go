package models

import (
	"time"

	"gorm.io/gorm"
)

type MailBox struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	MessageID string `gorm:"uniqueIndex;not null"`
	DateTS    int64  `gorm:"not null"` // Unix timestamp
	FromAddr  string `gorm:"index"`
	ToAddr    string `gorm:"index"`
	CcAddr    string
	BccAddr   string
	Subject   string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"` // optional, for soft delete
}
