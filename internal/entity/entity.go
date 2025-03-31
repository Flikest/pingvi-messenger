package entity

import (
	"time"

	"github.com/google/uuid"
)

type (
	User struct {
		ID       uuid.UUID `json:"id"`
		Name     string    `json:"name"`
		Pass     string    `json:"pass"`
		Email    string    `json:"email"`
		Avatar   string    `json:"avatar"`
		About_me string    `json:"about_me"`
	}

	Message struct {
		Chat_ID     string    `json:"chat_id"`
		Message_ID  int       `json:"message_id"`
		Sender_ID   string    `json:"sender_id"`
		Content     []byte    `json:"content"`
		SendingTime time.Time `json:"sending_time"`
	}

	Contact struct {
		User_id    uuid.UUID `json:"user_id"`
		Contact_id uuid.UUID `json:"contact_id"`
	}

	Participant struct {
		Chat_ID string `json:"chat_id"`
		User_ID string `json:"user_id"`
		IsAdmin bool   `json:"is_admin"`
	}

	Chat struct {
		ID              string    `json:"chat_id"`
		Name            string    `json:"chat_name"`
		Avatar          string    `json:"chat_avatar"`
		UniqueLinToJoin string    `json:"unique_link_to_join"`
		Last_seen       time.Time `json:last_seen`
	}

	Sticker struct {
	}
)
