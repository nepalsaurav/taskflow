package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func OpenDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	// language=sql
	_, err = db.Exec(`
		PRAGMA busy_timeout       = 10000;
		PRAGMA journal_mode       = WAL;
		PRAGMA journal_size_limit = 200000000;
		PRAGMA synchronous        = NORMAL;
		PRAGMA foreign_keys       = ON;
		PRAGMA temp_store         = MEMORY;
		PRAGMA cache_size         = -32000;
	`)
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, err
}
