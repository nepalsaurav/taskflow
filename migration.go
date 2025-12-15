package main

import (
	"fmt"
	"log"
	"os"
	"taskflow/models"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

type Migrate struct{}

func (m Migrate) up() {
	m.createDir()
	goose.SetDialect("sqlite")

	dbPaths := map[string]string{
		"database/mail.db": "migrations/mail",
	}

	for dbPath, migrationDir := range dbPaths {
		db, err := models.DefaultDBConnect(dbPath)
		if err != nil {
			log.Fatal(err)
		}

		if err := goose.Up(db.DB(), migrationDir); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Database: %s, Migrations: %s\n", dbPath, migrationDir)
	}

}

func (m Migrate) createDir() {
	directories := []string{"database"}

	for _, dir := range directories {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Panicf("error creating directories %s", dir)
		}
	}
}
