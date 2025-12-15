package models

type MailBox struct {
	ID          int64  `db:"id" json:"id"`
	TrackingID  string `db:"tracking_id" json:"tracking_id"`
	MessageID   string `db:"message_id" json:"message_id"`
	MaildirPath string `db:"maildir_path" json:"maildir_path"`
	FromAddr    string `db:"from_addr" json:"from_addr"`
	ToAddr      string `db:"to_addr" json:"to_addr"`
	CCAddr      string `db:"cc_addr" json:"cc_addr,omitempty"`
	BCCAddr     string `db:"bcc_addr" json:"bcc_addr,omitempty"`
	Subject     string `db:"subject" json:"subject,omitempty"`
	DateTs      int64  `db:"date_ts" json:"date_ts"`
}
