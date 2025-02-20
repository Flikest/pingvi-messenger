package storage

import (
	"context"
	"database/sql"
)

type Storage struct {
	db      *sql.DB
	context context.Context
}

func NewStorage(db *sql.DB, context context.Context) *Storage {
	return &Storage{db: db, context: context}
}
