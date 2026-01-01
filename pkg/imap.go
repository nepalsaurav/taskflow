package pkg

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/nepalsaurav/taskflow/models"
	"gorm.io/gorm/clause"
)

const (
	GMAIL_SENTBOX = "[Gmail]/Sent Mail"
)

type ImapConfig struct {
	LastProcessSentBoxUID uint32 `json:"last_process_sentbox_uid"`
}

func getImapConfigPath() string {
	dir := filepath.Join(os.Getenv("HOME"), ".config", "taskflow")
	_ = os.MkdirAll(dir, 0755)
	return filepath.Join(dir, "imap_config.json")
}

func getImapConfig() ImapConfig {
	path := getImapConfigPath()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ImapConfig{}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return ImapConfig{}
	}
	var cfg ImapConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return ImapConfig{}
	}
	return cfg
}

func saveConfig(cfg ImapConfig) error {
	path := getImapConfigPath()
	data, _ := json.MarshalIndent(cfg, "", "  ")
	return os.WriteFile(path, data, 0600)
}

func connectImap() (*imapclient.Client, error) {
	client, err := imapclient.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		return nil, err
	}

	if err := client.Login("nepalsaurav123@gmail.com", "rdue nimu wdpb sqzu").Wait(); err != nil {
		return nil, err
	}

	return client, nil
}

func parseAddress(addresses []imap.Address) string {
	if len(addresses) == 0 {
		return ""
	}

	addrs := make([]string, 0, len(addresses))
	for _, a := range addresses {
		addrs = append(addrs, a.Addr())
	}

	return strings.Join(addrs, ",")
}

func ImapLogin() {

	client, err := connectImap()

	if err != nil {
		fmt.Println(err)
	}

	status, err := client.Status(GMAIL_SENTBOX, &imap.StatusOptions{UIDNext: true}).Wait()

	if err != nil {
		return
	}

	// check sent box exist or not
	if status.UIDNext <= 1 {
		// no email in sentbox
		return
	}

	currentMaxUID := status.UIDNext - 1

	// get imap config

	imapConfig := getImapConfig()

	if currentMaxUID <= imap.UID(imapConfig.LastProcessSentBoxUID) {
		// no new sent mail
		return
	}

	_, err = client.Select(GMAIL_SENTBOX, &imap.SelectOptions{ReadOnly: true}).Wait()

	if err != nil {
		return
	}

	targetRange := imap.SeqSet{}
	targetRange.AddRange(imapConfig.LastProcessSentBoxUID+1, uint32(currentMaxUID))

	fetchCommand := client.Fetch(targetRange, &imap.FetchOptions{
		Envelope: true,
		UID:      true,
	})

	mailboxes := []models.MailBox{}

	for {
		msg := fetchCommand.Next()
		if msg == nil {
			break
		}
		itemData, err := msg.Collect()
		if err != nil {
			continue
		}

		mailbox := models.MailBox{
			MessageID: itemData.Envelope.MessageID,
			DateTS:    itemData.Envelope.Date.Unix(),
			FromAddr:  parseAddress(itemData.Envelope.From),
			ToAddr:    parseAddress(itemData.Envelope.To),
			CcAddr:    parseAddress(itemData.Envelope.Cc),
			BccAddr:   parseAddress(itemData.Envelope.Bcc),
			Subject:   itemData.Envelope.Subject,
		}

		mailboxes = append(mailboxes, mailbox)

	}

	db, err := models.DefaultDBConnect("database/models.db")

	if err != nil {
		return
	}

	result := db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(mailboxes, 100)

	if result.Error == nil {
		imapConfig.LastProcessSentBoxUID = uint32(currentMaxUID)
		saveConfig(imapConfig)
	}

}
