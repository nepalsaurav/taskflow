package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
	"taskflow/db"

	sq "github.com/Masterminds/squirrel"
	"github.com/k0kubun/pp/v3"
)

// IndexMail scans the user's Maildir "new" directory and indexes all messages.
// It resolves the Maildir path, builds the path to "new/", and delegates the
// actual indexing work to IndexMailByPath. Returns an error only if the Maildir
// root cannot be resolved; errors during indexing are handled inside
// IndexMailByPath.
func IndexMail() error {
	maildirPath, err := GetMailDirPath()
	if err != nil {
		return err
	}
	maildirNewPath := filepath.Join(maildirPath, "new")
	IndexMailByPath(maildirNewPath)
	return nil
}

// IndexMailByPath walks the specified Maildir path and processes every message
// file it contains. Directories are skipped. For each file, it calls
// getEmailMessage to parse the email and logs or indexes the resulting
// MailMessage. File-level errors are logged and ignored so the walk continues.
// Returns nil unless the walk setup itself fails.
func IndexMailByPath(path string) error {
	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("unable to read maildir path: %v\n", err)
			return nil
		}
		if d.IsDir() {
			return nil
		}

		emailMessage, _ := getEmailMessage(path)

		db, err := db.OpenDB("database/mailbox.db")
		if err != nil {
			println(err)
		}
		defer db.Close()

		query := sq.Insert("mailbox").
			Columns("tracking_id", "message_id", "maildir_path", "from_addr", "to_addr", "cc_addr", "bcc_addr", "subject", "body_text", "date_ts", "status").
			Values(
				emailMessage.TrackingID,
				emailMessage.MessageID,
				emailMessage.MaildirPath,
				emailMessage.FromAddr,
				emailMessage.ToAddr,
				emailMessage.CCAddr,
				emailMessage.BCCAddr,
				emailMessage.Subject,
				emailMessage.BodyText,
				emailMessage.DateTS,
				emailMessage.Status,
			).PlaceholderFormat(sq.Question)

		sqlStr, args, err := query.ToSql()

		if err != nil {
			log.Printf("failed to build sql: %v\n", err)
			return nil
		}

		result, err := db.Exec(sqlStr, args...)

		if err != nil {
			log.Printf("failed to insert mailbox: %v\n", err)
			return nil
		}

		fmt.Println(result)
		pp.Println(emailMessage)

		return nil
	})

	return nil
}

// getEmailMessage reads and parses an email file at the given path into a MailMessage struct.
// It extracts common headers (From, To, Cc, Bcc, Subject, Message-ID) and the message body.mail
// The MailDate field is set if the "Date" header can be parsed. FileId is derived from the
// Maildir filename using GetUniqueFileId. Any file-level or parsing errors are logged;mail
// in case of error, an empty MailMessage is returned along with nil error to allow processing
// to continue. Status is set to "Sent" by default.
func getEmailMessage(path string) (db.Mailbox, error) {
	fileData, err := os.ReadFile(path)
	if err != nil {
		log.Printf("failed to read file %s: %v", path, err)
		return db.Mailbox{}, nil
	}

	msg, err := mail.ReadMessage(bytes.NewReader(fileData))
	if err != nil {
		log.Printf("failed to parse email %s: %v", path, err)
		return db.Mailbox{}, nil
	}

	body, err := io.ReadAll(msg.Body)
	msgBody := ""
	if err != nil {
		log.Printf("failed to read body %s: %v", path, err)
	}
	msgBody = string(body)

	var mailMessage db.Mailbox
	mailMessage.FromAddr = msg.Header.Get("From")
	mailMessage.ToAddr = parseAddressList(msg.Header.Get("To"))
	mailMessage.CCAddr = parseAddressList(msg.Header.Get("Cc"))
	mailMessage.BCCAddr = parseAddressList(msg.Header.Get("Bcc"))
	mailMessage.Subject = msg.Header.Get("Subject")
	mailMessage.BodyText = msgBody
	mailMessage.Status = "Sent"
	mailMessage.MaildirPath = path
	mailMessage.MessageID = msg.Header.Get("Message-ID")
	mailMessage.TrackingID = msg.Header.Get("Tracking-ID")

	if date, err := msg.Header.Date(); err == nil {
		mailMessage.DateTS = date.Unix()
	}

	return mailMessage, nil
}

// get maildir path
func GetMailDirPath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("unable to get homedir path err: %v\n", err)
		return string(""), fmt.Errorf("unable to get homedir path err: %v\n", err)
	}
	maildirPath := filepath.Join(homedir, "Maildir")
	return maildirPath, nil
}

func parseAddressList(header string) string {
	addrs, err := mail.ParseAddressList(header)
	if err != nil {
		return ""
	}

	list := make([]string, len(addrs))
	for i, a := range addrs {
		list[i] = a.Address
	}
	return strings.Join(list, "")
}
