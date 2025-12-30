package models

import (
	"time"

	"gorm.io/gorm"
)

type MailBox struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	TrackingID  string `gorm:"uniqueIndex;not null"`
	MessageID   string `gorm:"uniqueIndex;not null"`
	MaildirPath string `gorm:"uniqueIndex;not null"`
	DateTS      int64  `gorm:"not null"` // Unix timestamp
	FromAddr    string `gorm:"index"`
	ToAddr      string `gorm:"index"`
	CcAddr      string
	BccAddr     string
	Subject     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"` // optional, for soft delete
}
