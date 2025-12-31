package models

import (
	"time"

	"gorm.io/gorm"
)

// MailLog represents a Postfix log entry
type MailLog struct {
	ID         uint          `gorm:"primaryKey"` // auto increment
	UniqueID   string        `gorm:"uniqueIndex;not null"`
	Date       time.Time     `gorm:"not null"`
	Status     string        `gorm:"size:50"`
	Receipents string        `gorm:"size:255"`
	MessageID  string        `gorm:"size:255"`
	Logs       []MailLogLine `gorm:"foreignKey:MailLogID;constraint:OnDelete:CASCADE"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

// MailLogLine represents a single line in the log slice
type MailLogLine struct {
	ID        uint   `gorm:"primaryKey"`
	MailLogID uint   `gorm:"index"`     // foreign key
	Line      string `gorm:"type:text"` // full log line
	CreatedAt time.Time
	UpdatedAt time.Time
}
