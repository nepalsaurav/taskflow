package pkg

import (
	"fmt"
	"io/fs"
	"log"
	"net/mail"
	"os"
	"path/filepath"
	"strings"

	"github.com/nepalsaurav/taskflow/models"
)

const (
	MAILDIR_DIR string = "Maildir"
	MAILDIR_NEW string = "new"
	MAILDIR_CUR string = "cur"
)

const (
	SQLITE_CONSTRAINT_UNIQUE = 2067
)

// MaildirConfig provides configuration utilities for locating Maildir paths.
type MaildirConfig struct{}

// Maildir handles indexing and processing of email messages stored in Maildir format.
type Maildir struct{}

// getDir constructs the full path to a subdirectory within the user's Maildir.
// Returns the path and any error encountered.
func (c MaildirConfig) getDir(path string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("error on getting user home dir, err: %v\n", err)
		return "", fmt.Errorf("error on getting user home dir, err: %w", err)
	}
	mailDir := filepath.Join(home, MAILDIR_DIR, path)
	return mailDir, nil
}

// getMailDirNew returns the path to the "new" subdirectory of the Maildir.
func (c MaildirConfig) getMailDirNew() (string, error) {
	return c.getDir(MAILDIR_NEW)
}

// getMailDirCur returns the path to the "cur" subdirectory of the Maildir.
func (c MaildirConfig) getMailDirCur() (string, error) {
	return c.getDir(MAILDIR_CUR)
}

// IndexMail scans the "new" Maildir directory and indexes any new emails into the database.
func (m Maildir) IndexMail() error {
	maildirConfig := MaildirConfig{}
	newPath, err := maildirConfig.getMailDirNew()
	if err != nil {
		return err
	}
	return m.indexMailByPath(newPath)

}

// indexMailByPath processes all files in the given directory and indexes them.
func (m Maildir) indexMailByPath(path string) error {
	return m.walkDir(path)
}

// walkDir walks the directory tree starting at path, parses email files concurrently,
// counts processed messages, and prints a summary of the operation.
func (m Maildir) walkDir(path string) error {

	err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}

		// parse email
		// msg, err := m.parseMail(p)
		// if err != nil {
		// 	return nil
		// }

		// // insert into database
		// if msg.TrackingID == "" {
		// 	msg.TrackingID = msg.MessageID
		// }

		return nil
	})
	// catch error of walkdir
	if err != nil {
		return fmt.Errorf("could not walk dir err: %w", err)
	}

	return nil

}

// parseMail reads an email file and extracts relevant metadata into a MailBox model.
func (m Maildir) parseMail(filePath string) (models.MailBox, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return models.MailBox{}, fmt.Errorf("can not load message file :%w", err)
	}
	defer file.Close()

	msg, err := mail.ReadMessage(file)

	mailMessage := models.MailBox{}

	// append data
	mailMessage.FromAddr = m.parseAddress(msg.Header.Get("From"))
	mailMessage.ToAddr = m.parseAddressList(msg.Header.Get("To"))
	mailMessage.CcAddr = m.parseAddressList(msg.Header.Get("CC"))
	mailMessage.BccAddr = m.parseAddressList(msg.Header.Get("BCC"))
	mailMessage.Subject = msg.Header.Get("Subject")
	mailMessage.MessageID = msg.Header.Get("Message-ID")

	if date, err := msg.Header.Date(); err == nil {
		mailMessage.DateTS = date.Unix()
	}
	return mailMessage, nil
}

// parseAddress extracts a single email address from a header string.
func (m Maildir) parseAddress(header string) string {
	addr, err := mail.ParseAddress(header)
	if err != nil {
		return ""
	}
	return addr.Address
}

// moveFile moves a file from the "new" directory to the "cur" directory.
// func (m Maildir) moveFile(filePath string) bool {
// 	maildirConfig := MaildirConfig{}
// 	curPath, err := maildirConfig.getMailDirCur()
// 	if err != nil {
// 		return false
// 	}

// 	newPath := filepath.Join(curPath, filepath.Base(filePath))

// 	if err := os.Rename(filePath, newPath); err != nil {
// 		return false
// 	}
// 	return true
// }

// parseAddressList extracts a comma-separated list of email addresses from a header.
func (m Maildir) parseAddressList(header string) string {
	addrs, err := mail.ParseAddressList(header)
	if err != nil {
		return ""
	}
	list := make([]string, len(addrs))
	for i, a := range addrs {
		list[i] = a.Address
	}
	return strings.Join(list, ",")
}
