package storage

import (
	"context"
	"database/sql"
	"sync"
)

type Storage struct {
	db      *sql.DB
	context context.Context
	mutex   sync.Mutex
}

func NewStorage(db *sql.DB, context context.Context) *Storage {
	return &Storage{db: db, context: context}
}
