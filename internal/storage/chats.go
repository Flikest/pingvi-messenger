package storage

import (
	"context"
	"database/sql"
	"fmt"
)

type Storage struct {
	db  *sql.DB
	ctx context.Context
}

func InitStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s Storage) AddMessages(mesages string) {
	query := fmt.Sprintf("INSERT INTO chats WHERE id=$1")
	s.db.ExecContext(s.ctx, query, mesages)
}

func (s Storage) GetChats()
