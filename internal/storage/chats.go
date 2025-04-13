package storage

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Flikest/PingviMessenger/internal/entity"
	"github.com/google/uuid"
)

func (s Storage) DataFromTheStartPage(user_ID string, ch chan []entity.Chat) {
	queryData := "SELECT * FROM chats JOIN messeges ON chats.id=messeges.chat_id WHERE messeges.sender_id = $1"

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

func (s Storage) GetAllMessageFromChat(chat_ID string, ch chan []entity.Message) {
	query := "SELECT * FROM messeges WHERE chat_id=$1"

	rows, err := s.db.QueryContext(s.context, query, chat_ID)
	if err != nil {
		slog.Info("sql query error", err)
	}

	var result = []entity.Message{}
	for rows.Next() {
		var message = entity.Message{}
		if err := rows.Scan(&message.Chat_ID, &message.Message_ID, &message.Sender_ID, &message.Content, &message.SendingTime); err != nil {
			slog.Info("can't read data from db", err)
		}
		result = append(result, message)
	}
	ch <- result
}

func (s Storage) CreateChat(creator_ID string, e entity.Chat) error {
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

func (s Storage) GetChat(chat_name string) ([]entity.Chat, error) {
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

	return result, err
}

func (s Storage) UpdateChat(body entity.Chat) (int, error) {
	query := "UPDATE chats SET name=$1, avatar=$2 WHERE id=$3"
	_, err := s.db.ExecContext(s.context, query, body.Name, body.Avatar, body.ID)
	if err != nil {
		slog.Info("failed to update data: ", err)
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

func (s Storage) AddMesage(message entity.Message) (int, error) {
	query := "INSERT INTO messeges (chat_id, message_id, sender_id, content, sending_time) VALUES ($1, $2, $3, $4, $5)"

	_, err := s.db.ExecContext(s.context, query, message.Chat_ID, message.Message_ID, message.Sender_ID, []byte(message.Content), message.SendingTime)
	if err != nil {
		slog.Info("error messages added: ", err)
		return http.StatusBadGateway, err
	}
	slog.Info("messages added")

	return http.StatusOK, nil
}

func (s Storage) GetMessage(chat_ID string, message_ID string) (entity.Message, error) {
	query := "SELECT mesegeges FROM chats where chat_id=$1 AND messege_id=$2"

	var message entity.Message

	row := s.db.QueryRowContext(s.context, query, chat_ID, message_ID).Scan(&message.Chat_ID, &message.Message_ID, &message.Sender_ID, &message.Content, &message.SendingTime)
	if row == nil {
		slog.Info("message not found!")
	}
	return message, nil
}

func (s Storage) UpdateMessage(e entity.Message) (int, error) {
	query := "UPDATE messeges SET content=$1 where chat_ID=$2 mesasge_ID=$3"

	_, err := s.db.ExecContext(s.context, query, []byte(e.Content), e.Chat_ID, e.Message_ID)
	if err != nil {
		slog.Info("failed to update message")
		return 404, err
	}
	return 200, nil
}

func (s Storage) DelelteMessage(chat_ID string, message_ID int) (int, error) {
	query := "DELETE messeges from chats WHERE message_id = $1 AND chat_id = $2"

	_, err := s.db.ExecContext(s.context, query, message_ID)
	if err != nil {
		slog.Info("failed to delete message")
		return http.StatusBadGateway, err
	}
	return http.StatusOK, nil
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
