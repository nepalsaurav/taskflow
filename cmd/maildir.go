package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"taskflow/models"
	"time"

	"github.com/k0kubun/pp/v3"
	"github.com/pocketbase/dbx"
)

const (
	MAILDIR_DIR string = "Maildir"
	MAILDIR_NEW string = "cur"
	MAILDIR_CUR string = "new"
)

type MaildirConfig struct{}
type Maildir struct{}
type IndexMailResp struct {
	numberOfMailIndex int
	Message           string
}

func (c MaildirConfig) getDir(path string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("error on getting user home dir, err: %v\n", err)
		return "", fmt.Errorf("error on getting user home dir, err: %w", err)
	}
	mailDir := filepath.Join(home, MAILDIR_DIR, path)
	return mailDir, nil
}

func (c MaildirConfig) getMailDirNew() (string, error) {
	return c.getDir(MAILDIR_NEW)
}

func (c MaildirConfig) getMailDirCur() (string, error) {
	return c.getDir(MAILDIR_CUR)
}

func (m Maildir) IndexMail() error {
	maildirConfig := MaildirConfig{}
	newPath, err := maildirConfig.getMailDirNew()
	if err != nil {
		return err
	}
	r, _ := m.indexMailByPath(newPath)
	pp.Println(r)
	return nil
}

func (m Maildir) indexMailByPath(path string) (IndexMailResp, error) {
	start := time.Now()
	mailMessageList, err := m.walkDir(path)
	if err != nil {
		return IndexMailResp{}, err
	}

	db, err := models.DefaultDBConnect("database/mail.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = db.Transactional(func(tx *dbx.Tx) error {
		for _, msg := range mailMessageList {
			if msg.TrackingID == "" {
				msg.TrackingID = msg.MessageID
			}

			_, err := tx.Insert("mailbox", dbx.Params{
				"tracking_id":  msg.TrackingID,
				"message_id":   msg.MessageID,
				"maildir_path": msg.MaildirPath,
				"date_ts":      msg.DateTs,
				"from_addr":    msg.FromAddr,
				"to_addr":      msg.ToAddr,
				"cc_addr":      msg.CCAddr,
				"bcc_addr":     msg.BCCAddr,
				"subject":      msg.Subject,
			}).Execute()

			if err != nil {
				return err
			}
			fmt.Println(msg.MessageID)
		}
		return nil
	})

	elapsed := time.Since(start)
	message := fmt.Sprintf("complete indexing mail in %s", elapsed)
	if len(mailMessageList) == 0 {
		message = "There is no new mail to index"
	}
	return IndexMailResp{numberOfMailIndex: len(mailMessageList), Message: message}, nil
}

func (m Maildir) walkDir(path string) ([]models.MailBox, error) {

	var wg sync.WaitGroup
	mailMessageList := make(chan models.MailBox)

	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("error %v", err)
		}
		if d.IsDir() {
			return nil
		}
		wg.Go(func() {
			msg, err := m.parseMail(path)
			if err != nil {
				return
			}
			mailMessageList <- msg
		})
		return nil
	})

	go func() {
		wg.Wait()
		close(mailMessageList)
	}()

	messageList := []models.MailBox{}
	for val := range mailMessageList {
		messageList = append(messageList, val)
	}

	return messageList, nil
}

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
	mailMessage.CCAddr = m.parseAddressList(msg.Header.Get("CC"))
	mailMessage.BCCAddr = m.parseAddressList(msg.Header.Get("BCC"))
	mailMessage.Subject = msg.Header.Get("Subject")
	mailMessage.MaildirPath = filePath
	mailMessage.MessageID = msg.Header.Get("Message-ID")
	mailMessage.TrackingID = msg.Header.Get("Tracking-ID")

	if date, err := msg.Header.Date(); err == nil {
		mailMessage.DateTs = date.Unix()
	}

	return mailMessage, nil
}

func (m Maildir) parseAddress(header string) string {
	addr, err := mail.ParseAddress(header)
	if err != nil {
		return ""
	}
	return addr.Address
}

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
