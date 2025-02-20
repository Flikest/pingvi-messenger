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
		About_me string    `json:"about_me"`
	}

	Messege struct {
		Message_id  int       `json:"message_id"`
		User_id     uuid.UUID `json:"User_id"`
		Content     []byte    `json:"content"`
		SendingTime time.Time `json:"sending_time"`
	}

	Sticker struct {
		Url_img string `json:"img"`
	}

	StickerPack struct {
		Name     string    `json:"name"`
		Сreator  string    `json:"name_creator"`
		Stickers []Sticker `json:"stickers"`
	}

	Participant struct {
		User_id  uuid.UUID `json:"user_id"`
		Name     string    `json:"name"`
		About_me string    `json:"about_me"`
	}

	Chat struct {
		ID         uuid.UUID `json:"chat_id"`
		Creator_id uuid.UUID `json:"creator_id"`
		Name       string    `json:"chat_name"`
		Img        string    `json:"chat_img"`
	}
)
