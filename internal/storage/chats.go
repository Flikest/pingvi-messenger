package storage

import (
	"fmt"
	"log/slog"

	"github.com/Flikest/PingviMessenger/internal/entity"
	"github.com/google/uuid"
)

func (s Storage) DataFromTheStartPage(user_ID string, ch chan []entity.Chat) {
	queryData := "SELECT DISTINCT * FROM chats messeges ON chats.id=messeges.chat_id WHERE messeges.sender_id = $1"

	rows, err := s.db.QueryContext(s.context, queryData, user_ID)
	if err != nil {
		slog.Info("sql query error", err)
	}

	var result = []entity.Chat{}

	for rows.Next() {
		chat := entity.Chat{}
		if err := rows.Scan(&chat.ID, &chat.Avatar, &chat.Name, &chat.Last_seen); err != nil {
			slog.Info("error retrieving home page data", err)
		}
		result = append(result, chat)
	}
	ch <- result
}

func (s Storage) CreateChat(creator_ID string, e entity.Chat, ch chan error) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	id := uuid.New()

	uniqueLink := fmt.Sprintf("http://pingui.org/messenger/?%s", id)

	queryChat := "INSERT INTO chats (id, name, avatar, unique_link_to_join) VALUES ($1, $2, $3, $4)"
	queryParticipants := "INSERT INTO participants (chat_id, user_id, is_admin) VALUES ($1, $2, $3)"

	_, err := s.db.ExecContext(s.context, queryChat, id, e.Name, e.Avatar, uniqueLink)
	if err != nil {
		slog.Info("failed to create chat: ", err)
		return err
	}

	_, err = s.db.ExecContext(s.context, queryParticipants, id, creator_ID, true)
	if err != nil {
		slog.Info("error adding creator: ", err)
		return err
	}

	return nil
}

func (s Storage) GetChat(ch chan []entity.Chat, chat_name string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	query := "SELECT * FROM chats WHERE name=$1"

	var chat entity.Chat

	rows, err := s.db.QueryContext(s.context, query, chat_name)

	var result []entity.Chat
	for rows.Next() {
		var row entity.Chat
		if err := rows.Scan(&chat.ID, &chat.Name, &chat.Avatar, &chat.UniqueLinToJoin, &chat.Last_seen); err != nil {
			slog.Info("Couldn't find columns: ", err)
		}
		result = append(result, row)
	}

	if err != nil {
		slog.Info("chat not found: ", err)
	}

	ch <- result
}

func (s Storage) UpdateChat(body entity.Chat, ch chan error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	query := "UPDATE chats SET name=$1, avatar=$2 WHERE id=$3"
	_, err := s.db.ExecContext(s.context, query, body.Name, body.Avatar, body.ID)
	if err != nil {
		slog.Info("failed to update data: ", err)
		ch <- err
	}
}

func (s Storage) DeleteChat(chat_ID string, ch chan error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	query := "DELETE FROM chats WHERE id=$1"

	_, err := s.db.ExecContext(s.context, query, chat_ID)
	if err != nil {
		slog.Info("Failed to delete user")
		ch <- err
	}
}

func (s Storage) AddUser(chat_ID string, user_ID string, ch chan error) {
	query := "INSERT INTO participants (chat_id, user_id, is_admin) VALUES ($1, $2, $3)"

	_, err := s.db.ExecContext(s.context, query, user_ID, chat_ID)
	if err != nil {
		slog.Info("failed to add user: ", err)
		ch <- err
	}
}

func (s Storage) DropUserFromChat(user_ID string, chat_ID string, ch chan error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	query := "DELETE FROM participants WHERE chat_id=$1 AND user_id=$2"

	_, err := s.db.ExecContext(s.context, query, user_ID, chat_ID)
	if err != nil {
		slog.Info("failed to kick user: ", err)
		ch <- err
	}
}
