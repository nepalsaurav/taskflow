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
	"taskflow/models"

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

		pp.Println(emailMessage)

		return nil
	})

	return nil
}

// getEmailMessage reads and parses an email file at the given path into a MailMessage struct.
// It extracts common headers (From, To, Cc, Bcc, Subject, Message-ID) and the message body.
// The MailDate field is set if the "Date" header can be parsed. FileId is derived from the
// Maildir filename using GetUniqueFileId. Any file-level or parsing errors are logged;
// in case of error, an empty MailMessage is returned along with nil error to allow processing
// to continue. Status is set to "Sent" by default.
func getEmailMessage(path string) (models.MailMessage, error) {
	fileData, err := os.ReadFile(path)
	if err != nil {
		log.Printf("failed to read file %s: %v", path, err)
		return models.MailMessage{}, nil
	}

	msg, err := mail.ReadMessage(bytes.NewReader(fileData))
	if err != nil {
		log.Printf("failed to parse email %s: %v", path, err)
		return models.MailMessage{}, nil
	}

	body, err := io.ReadAll(msg.Body)
	msgBody := ""
	if err != nil {
		log.Printf("failed to read body %s: %v", path, err)
	}
	msgBody = string(body)

	var mailMessage models.MailMessage
	mailMessage.From = msg.Header.Get("From")
	mailMessage.To = parseAddressList(msg.Header.Get("To"))
	mailMessage.CC = parseAddressList(msg.Header.Get("Cc"))
	mailMessage.BCC = parseAddressList(msg.Header.Get("Bcc"))
	mailMessage.Subject = msg.Header.Get("Subject")
	mailMessage.Body = msgBody
	mailMessage.Status = "Sent"
	mailMessage.FileId = GetUniqueFileId(path)
	mailMessage.MessageID = msg.Header.Get("Message-ID")

	if date, err := msg.Header.Date(); err == nil {
		mailMessage.MailDate = date
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

func parseAddressList(header string) []string {
	addrs, err := mail.ParseAddressList(header)
	if err != nil {
		return []string{}
	}

	list := make([]string, len(addrs))
	for i, a := range addrs {
		list[i] = a.Address
	}
	return list
}

func GetUniqueFileId(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) >= 2 {
		return parts[1]
	}
	return filename
}
