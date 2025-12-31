package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func DefaultDBConnect(dbPath string) (*gorm.DB, error) {
	pragmas := "?_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)&_pragma=journal_size_limit(200000000)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(ON)&_pragma=temp_store(MEMORY)&_pragma=cache_size(-32000)"

	db, err := gorm.Open(sqlite.Open(dbPath+pragmas), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
