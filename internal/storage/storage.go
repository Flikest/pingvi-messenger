package storage

import "database/sql"

type Storage struct {
	db *sql.DB
}

func InitStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}
