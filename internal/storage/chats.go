package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Flikest/PingviMessenger/internal/entity"
	"github.com/google/uuid"
)

type Storage struct {
	db  *sql.DB
	ctx context.Context
}

func InitStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s Storage) AddMessages(mesages entity.Msg, id uuid.UUID) {
	query := fmt.Sprintf("UPDATE chats SET messages = $1 WHERE id = $2")
	s.db.ExecContext(s.ctx, query, mesages, id)
}

func (s Storage) GetChats(chatName string) {
	query := fmt.Sprintf("SELECT * FROM chats WHERE chat_name=$1")
	s.db.QueryContext(s.ctx, query, chatName)
}

func (s Storage) UpdateChats(c entity.Chats) {
	query := fmt.Sprintf("UPDATE chats SET chat_name = $1, img = $2, is_publick = $3 WHERE id = $4")
	s.db.QueryContext(s.ctx, query, c.Chat_name, c.Img, c.Is_publisk, c.ID)
}

func (s Storage) DeleteChat(id uuid.UUID) {
	query := fmt.Sprintf("DELETE FROM chats WHERE id = $1")
	s.db.QueryContext(s.ctx, query, id)
}
