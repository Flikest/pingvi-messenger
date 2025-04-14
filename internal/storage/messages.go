package storage

import (
	"log/slog"

	"github.com/Flikest/PingviMessenger/internal/entity"
)

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

func (s Storage) AddMesage(ch chan error, message entity.Message) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	query := "INSERT INTO messeges (chat_id, message_id, sender_id, content, sending_time) VALUES ($1, $2, $3, $4, $5)"

	_, err := s.db.ExecContext(s.context, query, message.Chat_ID, message.Message_ID, message.Sender_ID, []byte(message.Content), message.SendingTime)
	if err != nil {
		slog.Info("error messages added: ", err)
		ch <- err
	}
	slog.Info("messages added")
}

func (s Storage) GetMessage(chat_ID string, message_ID string, ch chan entity.Message) {
	query := "SELECT mesegeges FROM chats where chat_id=$1 AND messege_id=$2"

	var message entity.Message

	row := s.db.QueryRowContext(s.context, query, chat_ID, message_ID).Scan(&message.Chat_ID, &message.Message_ID, &message.Sender_ID, &message.Content, &message.SendingTime)
	if row == nil {
		slog.Info("message not found!")
	}
	ch <- message
}

func (s Storage) UpdateMessage(ch chan error, e entity.Message) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	query := "UPDATE messeges SET content=$1 where chat_ID=$2 mesasge_ID=$3"

	_, err := s.db.ExecContext(s.context, query, []byte(e.Content), e.Chat_ID, e.Message_ID)
	if err != nil {
		slog.Info("failed to update message")
		ch <- err
	}
}

func (s Storage) DeleteMessage(ch chan error, chat_ID string, message_ID int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	query := "DELETE messeges from chats WHERE message_id = $1 AND chat_id = $2"

	_, err := s.db.ExecContext(s.context, query, message_ID)
	if err != nil {
		slog.Info("failed to delete message")
		ch <- err
	}
}
