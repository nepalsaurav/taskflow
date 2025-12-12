package models

import "time"

type MailMessage struct {
	ID        uint
	MailDate  time.Time
	FileId    string
	MessageID string
	From      string
	To        []string
	CC        []string
	BCC       []string
	Subject   string
	Body      string
	Status    string
	CreatedAt *time.Time
}
