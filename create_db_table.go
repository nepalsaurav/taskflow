package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"taskflow/db"
	"time"
)

// Define each DB and its SQL file
var dbConfigs = []struct {
	Name    string
	SQLFile string
}{
	{"mailbox.db", "db_table/mailbox.sql"},
}

const (
	dbDir     = "database"
	backupDir = "backup_database"
)

func CreateDBTable() {
	// Ensure directories exist
	ensureDir(dbDir)
	ensureDir(backupDir)

	// Iterate over all DB configs
	for _, cfg := range dbConfigs {
		dbPath := filepath.Join(dbDir, cfg.Name)
		log.Printf("Processing DB: %s", dbPath)

		// Backup and delete old DB
		backupAndDeleteOldDB(dbPath)

		// Create new DB safely
		createDBTableSafe(dbPath, cfg.SQLFile)
	}
}

func ensureDir(dir string) {
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = sourceFile.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = destFile.ReadFrom(sourceFile)
	return err
}

func backupAndDeleteOldDB(dbPath string) {
	if _, err := os.Stat(dbPath); err == nil {
		timestamp := time.Now().Format("20060102_150405")
		backupPath := filepath.Join(backupDir, fmt.Sprintf("%s_%s", filepath.Base(dbPath), timestamp))
		if err := copyFile(dbPath, backupPath); err != nil {
			log.Fatalf("Failed to backup old DB %s: %v", dbPath, err)
		}
		log.Printf("Old DB backed up to %s", backupPath)

		if err := os.Remove(dbPath); err != nil {
			log.Fatalf("Failed to delete old DB %s: %v", dbPath, err)
		}
		log.Printf("Old DB %s deleted", dbPath)
	}
}

func createDBTableSafe(dbPath, sqlFile string) {
	// Read SQL content
	sqlContent, err := os.ReadFile(sqlFile)
	if err != nil {
		log.Panicf("Error reading SQL file %s: %v", sqlFile, err)
	}

	tmpDB := dbPath + ".tmp"

	// Remove temp DB if exists
	if _, err := os.Stat(tmpDB); err == nil {
		os.Remove(tmpDB)
	}

	dbConn, err := db.OpenDB(tmpDB)
	if err != nil {
		log.Panicf("Error opening temp DB %s: %v", tmpDB, err)
	}
	defer dbConn.Close()

	// Execute SQL
	_, err = dbConn.Exec(string(sqlContent))
	if err != nil {
		log.Panicf("Error creating tables in temp DB %s: %v", tmpDB, err)
	}

	log.Printf("Temporary DB %s created successfully", tmpDB)

	// Rename temp DB to final DB
	if err := os.Rename(tmpDB, dbPath); err != nil {
		log.Fatalf("Failed to rename temp DB to final DB %s: %v", dbPath, err)
	}

	log.Printf("New DB %s created successfully", dbPath)
}
