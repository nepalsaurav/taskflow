package pkg

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/nepalsaurav/taskflow/models"
	"gorm.io/gorm/clause"
)

var (
	reQueueID      = regexp.MustCompile(`\b([0-9A-F]{8,16})\b`)
	reTo           = regexp.MustCompile(`to=<([^>]+)>`)
	reMessageID    = regexp.MustCompile(`message-id=<([^>]+)>`)
	reStatus       = regexp.MustCompile(`status=(\S+)`)
	PostfixLogPath = "/var/log/mail.log"
	BatchSize      = 500
)

// analyzeLine parses a single line and updates the log entry in the map.
func analyzeLine(line string, logDict map[string]*models.MailLog) {
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return
	}

	// 1. Extract Queue ID
	match := reQueueID.FindStringSubmatch(line)
	if len(match) < 2 {
		return
	}
	uniqueID := match[1]

	// 2. Get or Create Entry
	entry, exists := logDict[uniqueID]
	if !exists {
		entry = &models.MailLog{
			UniqueID: uniqueID,
		}
		logDict[uniqueID] = entry
	}

	parts := strings.SplitN(line, " ", 2)
	if len(parts) > 0 {
		if parsedDate, err := time.Parse(time.RFC3339, parts[0]); err == nil {
			entry.Date = parsedDate
		}
	}

	if m := reTo.FindStringSubmatch(line); len(m) > 1 {
		entry.Receipents += "," + m[1]
	}

	if m := reMessageID.FindStringSubmatch(line); len(m) > 1 {
		entry.MessageID = m[1]
	}

	if m := reStatus.FindStringSubmatch(line); len(m) > 1 {
		entry.Status = m[1]
	}

	// Append Raw Log Line
	entry.Logs = append(entry.Logs, models.MailLogLine{
		Line: line,
	})
}

// saveBatch handles the database insertion logic
func saveBatch(logDict map[string]*models.MailLog) error {
	// Connect to DB once
	db, err := models.DefaultDBConnect("database/models.db")
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}
	var logs []*models.MailLog
	for _, logEntry := range logDict {
		logs = append(logs, logEntry)
	}
	if len(logs) == 0 {
		return nil
	}

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "unique_id"}},
		UpdateAll: true,
	}).Create(&logs)

	return result.Error
}

func ScanLogFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("unable to open log file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// Increase buffer size if lines are very long
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	logDict := make(map[string]*models.MailLog)

	for scanner.Scan() {
		analyzeLine(scanner.Text(), logDict)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	saveBatch(logDict)

	return nil
}

func GetMailLog() {
	if err := ScanLogFile(PostfixLogPath); err != nil {
		log.Printf("Failed to scan logs: %v", err)
	}
}
