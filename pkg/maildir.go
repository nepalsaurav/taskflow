package pkg

import (
	"fmt"
	"io/fs"
	"log"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/k0kubun/pp"
	"github.com/nepalsaurav/postfix_admin/models"
	"github.com/pocketbase/dbx"
	"modernc.org/sqlite"
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

// IndexMailResp represents the result of indexing emails.
type IndexMailResp struct {
	numberOfMailIndex int
	Message           string
}

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
	m.indexMailByPath(newPath)
	return nil
}

// indexMailByPath processes all files in the given directory and indexes them.
func (m Maildir) indexMailByPath(path string) (IndexMailResp, error) {
	m.walkDir(path)
	return IndexMailResp{}, nil
}

// walkDir walks the directory tree starting at path, parses email files concurrently,
// counts processed messages, and prints a summary of the operation.
func (m Maildir) walkDir(path string) (IndexMailResp, error) {
	start := time.Now()

	// counter for number of files scanned
	var counter int64

	// open database
	db, err := models.DefaultDBConnect("database/postfix_admin.db")
	if err != nil {
		fmt.Println(err)
		return IndexMailResp{}, err
	}

	// start transaction
	tx, err := db.Begin()
	if err != nil {
		return IndexMailResp{}, err
	}

	// record succesfull parse and record file
	successFile := []string{}
	// walk directory
	err = filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}

		// parse email
		msg, err := m.parseMail(p)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		// insert into database
		if msg.TrackingID == "" {
			msg.TrackingID = msg.MessageID
		}
		_, err = tx.Insert("mailbox", dbx.Params{
			"tracking_id":  msg.TrackingID,
			"message_id":   msg.MessageID,
			"maildir_path": msg.MaildirPath,
			"date_ts":      msg.DateTS,
			"from_addr":    msg.FromAddr,
			"to_addr":      msg.ToAddr,
			"cc_addr":      msg.CcAddr,
			"bcc_addr":     msg.BccAddr,
			"subject":      msg.Subject,
		}).Execute()

		// check sql error
		if err != nil {
			if sqlerr, ok := err.(*sqlite.Error); ok {
				if sqlerr.Code() == SQLITE_CONSTRAINT_UNIQUE {
					successFile = append(successFile, msg.MaildirPath)
				}
			}
		}

		if err == nil {
			successFile = append(successFile, msg.MaildirPath)
		}

		counter++
		return nil
	})
	// catch error of walkdir
	if err != nil {
		return IndexMailResp{}, fmt.Errorf("could not walk dir err: %w", err)
	}

	// commit transaction
	_ = tx.Commit()

	// check success file and move file from unread to read
	for _, file := range successFile {
		m.moveFile(file)
	}

	// response message
	elapse := time.Since(start)
	message := fmt.Sprintf("%d total email message scan and sync to database in %s", counter, elapse)
	if counter == 0 {
		message = "no new email to scan"
	}

	responseMessage := IndexMailResp{
		numberOfMailIndex: int(counter),
		Message:           message,
	}

	pp.Println(responseMessage)
	return responseMessage, nil
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
	mailMessage.MaildirPath = filePath
	mailMessage.MessageID = msg.Header.Get("Message-ID")
	mailMessage.TrackingID = msg.Header.Get("Tracking-ID")

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
func (m Maildir) moveFile(filePath string) bool {
	maildirConfig := MaildirConfig{}
	curPath, err := maildirConfig.getMailDirCur()
	if err != nil {
		return false
	}

	newPath := filepath.Join(curPath, filepath.Base(filePath))

	if err := os.Rename(filePath, newPath); err != nil {
		return false
	}
	return true
}

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

func (m Maildir) WatchMailDir(watcher *fsnotify.Watcher) {
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Create) {
					// parse mail
					parsedMail, err := m.parseMail(event.Name)
					if err != nil {
						log.Printf("error while parsing mail from watch maildir err: %v\n", err)
					}

					_ = m.moveFile(event.Name)

					pp.Println(parsedMail)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	maildirConfig := MaildirConfig{}
	newPath, err := maildirConfig.getMailDirNew()
	if err != nil {
		log.Printf("unable to get maildir new path, err: %v\n", err)
	}

	fmt.Println(newPath)

	err = watcher.Add(newPath)
	if err != nil {
		log.Printf("unable to watch maildir new path, err: %v\n", err)
	}

}
