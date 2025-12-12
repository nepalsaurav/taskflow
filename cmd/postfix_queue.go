package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/k0kubun/pp/v3"
)

type Recipient struct {
	Address     string `json:"address"`
	DelayReason string `json:"delay_reason"`
}

type QueueEntry struct {
	QueueName    string      `json:"queue_name"`
	QueueID      string      `json:"queue_id"`
	ArrivalTime  int64       `json:"arrival_time"`
	MessageSize  int         `json:"message_size"`
	ForcedExpire bool        `json:"forced_expire"`
	Sender       string      `json:"sender"`
	Recipients   []Recipient `json:"recipients"`
}

func GetPostfixQueue() ([]QueueEntry, error) {
	cmd := exec.Command("postqueue", "-j")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	var entries []QueueEntry
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		var entry QueueEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			fmt.Println("Failed to parse line:", line, "error:", err)
			continue
		}
		entries = append(entries, entry)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	pp.Println(entries)

	return entries, nil
}
