package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Flikest/PingviMessenger/internal/entity"
	"github.com/google/uuid"
)

func (s Storage) CreateChat(creator_ID string, e entity.Chat) (uuid.UUID, error) {
	id := uuid.New()

	uniqueLink := fmt.Sprintf("http://pingui.org/messenger/?%s", id)

	queryChat := "INSERT INTO chats (id, name, avatar, unique_link_to_join) VALUES ($1, $2, $3, $4)"
	queryParticipants := "INSERT INTO participants (chat_id, user_id, is_admin) VALUES ($1, $2, $3)"

	_, err := s.db.ExecContext(s.context, queryChat, id, e.Name, e.Avatar, uniqueLink)
	if err != nil {
		slog.Info("failed to create chat: ", err)
		return uuid.Nil, err
	}

	_, err = s.db.ExecContext(s.context, queryParticipants, id, creator_ID, true)
	if err != nil {
		slog.Info("error adding creator: ", err)
		return uuid.Nil, err
	}

	return id, nil
}

func (s Storage) GetChat(chat_ID string) (*sql.Rows, error) {
	query := "SELECT * FROM chats WHERE id=$1"

	chat, err := s.db.QueryContext(s.context, query, chat_ID)
	if err != nil {
		slog.Info("chat not found")
		return nil, errors.New("nothing found!")
	}
	return chat, nil
}

func (s Storage) UpdateChat(body entity.Chat) (int, error) {
	query := "UPDATE chats SET name=$1, avatar=$2 WHERE id=$3"
	_, err := s.db.ExecContext(s.context, query, body.Name, body.Avatar, body.ID)
	if err != nil {
		slog.Info("failed to update data", err)
		return 404, err
	}
	return 200, nil
}

func (s Storage) DeleteChat(chat_ID string) (int, error) {
	query := "DELETE FROM chats WHERE id=$1"

	_, err := s.db.ExecContext(s.context, query, chat_ID)
	if err != nil {
		slog.Info("Failed to delete user")
		return 404, err
	}
	return 200, nil
}

func (s Storage) AddMesage(message entity.Messege) (int, error) {
	query := "INSERT INTO messeges (chat_id, message_id, sender_id, content, sending_time) VALUES ($1, $2, $3, $4, $5)"

	_, err := s.db.ExecContext(s.context, query, message.Chat_ID, message.Message_ID, message.Sender_ID, []byte(message.Content), message.SendingTime)
	if err != nil {
		slog.Info("error messages added: ", err)
		return 404, err
	}
	slog.Info("messages added")
	return 200, nil
}

func (s Storage) GetMessage(chat_ID string, message_ID string) (*sql.Row, error) {
	query := "SELECT mesegeges FROM chats where chat_id=$1 AND messege_id=$2"

	message := s.db.QueryRowContext(s.context, query, chat_ID, message_ID)
	if message == nil {
		return nil, errors.New("message not found!")
	}
	return message, nil
}

func (s Storage) UpdateMessage(e entity.Messege) (int, error) {
	query := "UPDATE messeges SET content=$1 where chat_ID=$2 mesasge_ID=$3"

	_, err := s.db.ExecContext(s.context, query, []byte(e.Content), e.Chat_ID, e.Message_ID)
	if err != nil {
		slog.Info("failed to update message")
		return 404, err
	}
	return 200, nil
}

func (s Storage) DelelteMessage(message_ID string) (int, error) {
	query := "DELETE messeges from chats WHERE message_ID = $1"

	_, err := s.db.ExecContext(s.context, query, message_ID)
	if err != nil {
		slog.Info("failed to delete message")
		return 404, err
	}
	return 200, nil
}

func (s Storage) AddUser(chat_ID string, user_ID string) (int, error) {
	query := "INSERT INTO participants (chat_id, user_id, is_admin) VALUES ($1, $2, $3)"

	_, err := s.db.ExecContext(s.context, query, user_ID, chat_ID)
	if err != nil {
		slog.Info("failed to add user: ", err)
		return 404, err
	}
	return 200, nil
}

func (s Storage) DropUserFromChat(user_ID string, chat_ID string) (int, error) {
	query := "DELETE FROM participants WHERE chat_id=$1 AND user_id=$2"

	_, err := s.db.ExecContext(s.context, query, user_ID, chat_ID)
	if err != nil {
		slog.Info("failed to kick user: ", err)
		return 404, err
	}
	return 200, nil
}

func (s Storage) CountMessages(chat_ID string) (int, error) {
	query := "SELECT COUNT(*) FROM messeges WHERE chat_id=$1"

	count := s.db.QueryRowContext(s.context, query, chat_ID)
	if count == nil {
		return 0, errors.New("nothing found")
	}

	var quantity int
	err := count.Scan(&quantity)
	if err != nil {
		return 0, err
	}

	return quantity, nil
}
