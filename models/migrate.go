package models

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	SQLITE_PRAGMA = "?_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)&_pragma=journal_size_limit(200000000)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(ON)&_pragma=temp_store(MEMORY)&_pragma=cache_size(-32000)"
)

func Migrate() {

	if err := os.MkdirAll("database", 0755); err != nil {
		log.Fatal("failed to create database folder:", err)
	}

	db, err := gorm.Open(sqlite.Open("database/models.db"+SQLITE_PRAGMA), &gorm.Config{})

	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	if err := db.AutoMigrate(&MailLog{}, &MailLogLine{}, &MailBox{}); err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	log.Println("Migration completed successfully")

}
