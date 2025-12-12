package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"log"
)

// Regex constants
var (
	TimestampRegex = regexp.MustCompile(`^([0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}(\.[0-9]+)?([Zz]|([+-][0-9]{2}:[0-9]{2})?))`)
	queueIDRegex   = regexp.MustCompile(`([A-F0-9]{10}):`)
	toRegex        = regexp.MustCompile(`to=<([^>@]+@[^>@]+)>`)
	fromRegex      = regexp.MustCompile(`from=<([^>@]+@[^>@]+)>`)
	messageIDRegex = regexp.MustCompile(`message-id=<(.*?)>`)
	statusRegex    = regexp.MustCompile(`status=([a-zA-Z0-9-]+)(?: (.*))?`)
	relayRegex     = regexp.MustCompile(`relay=([^\[]+)\[([^\]]+)\]:(\d+)`)
	clientRegex    = regexp.MustCompile(`client=([^\[]+)\[([^\]]+)\]`)
)

type PostfixLogEntry struct {
	Timestamp  string
	QueueID    string
	Message    string
	To         string
	From       string
	MessageID  string
	Status     string
	StatusDesc string
	RelayHost  string
	RelayIP    string
	RelayPort  string
	ClientHost string
	ClientIP   string
}

func ParsePostfixLog() error {
	cfg, err := GetSystemConfig()
	if err != nil {
		return err
	}

	logFile, err := os.Open(cfg.PostfixLogFile)
	if err != nil {
		log.Printf("error while opening postfix log file, err: %v\n", err)
		return fmt.Errorf("error while opening postfix log file, err: %v\n", err)
	}
	defer logFile.Close()

	scanner := bufio.NewScanner(logFile)

	for scanner.Scan() {

		var logEntry PostfixLogEntry

		line := scanner.Text()
		if line == "" {
			continue
		}

		if timestamp := TimestampRegex.FindStringSubmatch(line); len(timestamp) > 1 {
			logEntry.Timestamp = timestamp[1]
		}
		if queueID := queueIDRegex.FindStringSubmatch(line); len(queueID) > 1 {
			logEntry.QueueID = queueID[1]
		}
		if to := toRegex.FindStringSubmatch(line); len(to) > 1 {
			logEntry.To = to[1]
		}
		if from := fromRegex.FindStringSubmatch(line); len(from) > 1 {
			logEntry.From = from[1]
		}

		if mid := messageIDRegex.FindStringSubmatch(line); len(mid) > 1 {
			logEntry.MessageID = mid[1]
		}

		if status := statusRegex.FindStringSubmatch(line); len(status) > 1 {
			logEntry.Status = status[1]
			if len(status) > 2 {
				logEntry.StatusDesc = status[2]
			}
		}
		if relay := relayRegex.FindStringSubmatch(line); len(relay) > 3 {
			logEntry.RelayHost = relay[1]
			logEntry.RelayIP = relay[2]
			logEntry.RelayPort = relay[3]
		}
		if client := clientRegex.FindStringSubmatch(line); len(client) > 2 {
			logEntry.ClientHost = client[1]
			logEntry.ClientIP = client[2]
		}
		// pp.Println(logEntry)

	}

	searchMaildirForMessageID()
	return nil
}

func searchMaildirForMessageID() error {

	return nil
}
