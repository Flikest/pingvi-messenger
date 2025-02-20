package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Flikest/PingviMessenger/internal/entity"
	"github.com/google/uuid"
)

func (s Storage) CreateChat(creator_id uuid.UUID) (int, error) {
	id := uuid.New()

	query := fmt.Sprintf("INSERT INTO Сhats (id, Participants) VALUES ($1, ARRAY_APPEND(Participants, $2)")

	_, err := s.db.ExecContext(s.context, query, id, creator_id)
	if err != nil {
		slog.Info("failed to create chat: ", err)
		return 404, err
	}
	return 200, nil
}

func (s Storage) GetChat(chat_name string) (*sql.Rows, error) {
	query := fmt.Sprintf("SELECT * FROM chats WHERE id=$1")

	chat, err := s.db.QueryContext(s.context, query, chat_name)
	if err != nil {
		slog.Info("nothing found")
		return nil, errors.New("nothing found!")
	}
	return chat, nil
}

func (s Storage) UpdateChat(body entity.Chat) (int, error) {
	query := fmt.Sprintf("UPDATE chats SET name=$1, img=$2 WHERE id=$3")
	_, err := s.db.ExecContext(s.context, query, body.Name, body.Img, body.ID)
	if err != nil {
		slog.Info("failed to update data", err)
		return 404, err
	}
	return 200, nil
}

func (s Storage) DeleteChat(chat_id string) (int, error) {
	query := fmt.Sprintf("DELETE FROM chats WHERE id=$1")

	_, err := s.db.ExecContext(s.context, query, chat_id)
	if err != nil {
		slog.Info("Failed to delete user")
		return 404, err
	}
	return 200, nil
}

func (s Storage) AddMesage(chat_id string, message entity.Messege) (int, error) {
	query := fmt.Sprintf("UPDATE chats SET messages=$1 WHERE id=$2")

	_, err := s.db.ExecContext(s.context, query, message, chat_id)
	if err != nil {
		slog.Info("error messages updated: ", err)
		return 404, err
	}
	slog.Info("messages updated")
	return 200, nil
}

func (s Storage) GetMessage(message_id string) (*sql.Row, error) {
	query := fmt.Sprintf("SELECT mesegeges FROM chats where messege_id=$1")

	message := s.db.QueryRowContext(s.context, query, message_id)
	if message != nil {
		return nil, errors.New("message not found!")
	}
	return message, nil
}

func (s Storage) UpdateMessage(message_id int, mesasge entity.Messege) (int, error) {
	query := fmt.Sprintf("UPDATE chats SET messages=$1 where id=$2")

	_, err := s.db.ExecContext(s.context, query)
	if err != nil {
		slog.Info("failed to update message")
		return 404, err
	}
	return 200, nil
}

func (s Storage) DelelteMessage(message_id string) (int, error) {
	query := fmt.Sprintf("DELETE mesages from chats WHERE id = $1")

	_, err := s.db.ExecContext(s.context, query, message_id)
	if err != nil {
		slog.Info("failed to delete message")
		return 404, err
	}
	return 200, nil
}

func (s Storage) AddUser(chat_id string, user_id uuid.UUID) (int, error) {
	query := fmt.Sprintf("UPDATE chats SET Participants=$1 WHERE id=$2")

	_, err := s.db.ExecContext(s.context, query, user_id, chat_id)
	if err != nil {
		slog.Info("failed to add user: ", err)
		return 404, err
	}
	return 200, nil
}

func (s Storage) DropUserFromChat(user entity.Participant, chat_id string) (int, error) {
	query := fmt.Sprintf("array_remove(chats, $1) FROM chats WHERE id=$2")

	_, err := s.db.ExecContext(s.context, query, user, chat_id)
	if err != nil {
		slog.Info("failed to kick user: ", err)
		return 404, err
	}
	return 200, nil
}
