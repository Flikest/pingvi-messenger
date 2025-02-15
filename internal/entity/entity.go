package entity

import (
	"time"

	"github.com/google/uuid"
)

type (
	Users struct {
		ID       uuid.UUID `json:"id"`
		Name     string    `json:"name"`
		Avatar   string    `json:"avatar"`
		About_me string    `json:"about_me"`
	}

	Msg struct {
		User_ID       uuid.UUID `json:"user_id"`
		Content       string    `json:"content"`
		Dispatch_time time.Time `json:"dispatch_time"`
		Dispatch_date time.Time `json:"Dispatch_date"`
		Reactions     []string  `json:"reactions"`
	}

	Post struct {
		Content       string    `json:"content"`
		Dispatch_time time.Time `json:"dispatch_time"`
		Dispatch_date time.Time `json:"dispatch_date"`
		Reactions     []Smiles  `json:""`
	}

	Smiles struct {
		Url_img string `json:"url_img"`
	}

	Sticker struct {
		Url_sticker string `json:"sticker"`
	}

	Stickerpack struct {
		ID uuid.UUID `json:"id"`
	}

	Chats struct {
		ID         uuid.UUID `json:"id"`
		Chat_name  string    `json:"chat_name"`
		Img        string    `json:"images"`
		Message    Msg       `json:"msg"`
		Is_publisk bool      `json:"is_publick"`
	}

	Groops struct {
		ID           uuid.UUID   `json:"ID"`
		Message      Msg         `json:"msg"`
		Participants []uuid.UUID `json:"participants"`
	}

	Channels struct {
		ID    uuid.UUID `json:"id"`
		Posts Post      `json:"posts"`
	}
)
