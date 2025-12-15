package models

import (
	"github.com/pocketbase/dbx"
	_ "modernc.org/sqlite"
)

func DefaultDBConnect(dbPath string) (*dbx.DB, error) {
	pragmas := "?_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)&_pragma=journal_size_limit(200000000)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(ON)&_pragma=temp_store(MEMORY)&_pragma=cache_size(-32000)"

	db, err := dbx.Open("sqlite", dbPath+pragmas)
	if err != nil {
		return nil, err
	}

	return db, nil
}
