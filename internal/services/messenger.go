package services

import (
	"log/slog"
	"net/http"

	"github.com/Flikest/PingviMessenger/internal/entity"
	"github.com/Flikest/PingviMessenger/internal/storage"

	"github.com/gorilla/websocket"
)

type Service struct {
	Storage *storage.Storage
}

var upgreader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: true,
}

func InitService(s *storage.Storage) *Service {
	return &Service{Storage: s}
}

func (s Service) Correspondence(w http.ResponseWriter, r *http.Request) {
	ws, err := upgreader.Upgrade(w, r, nil)
	if err != nil {
		slog.Debug("error during update: ", err)
	}
	defer ws.Close()

	if r.Header.Get("Upgrade") != "websocket" {
		return 
	}

	if r.Header.Get("Connection") != "Upgrade" {
		return
	}

	swk := r.Header.Get("Sec-WebSocket-Key") 
	if swk == "" {
		return
	}

	goud := 

	var body entity.Msg = entity.Msg{}
	for {
		err := ws.ReadJSON(&body)
		if err != nil {
			slog.Info("reading error: ", err)
			return
		}
		if err := ws.WriteJSON(&body); err != nil {
			slog.Info("writing error: ", err)
			return
		}
		w.Write([]byte("привет"))
		s.Storage.AddMessages(body, body.User_ID)
	}
}

func (s Service) UpdateMessage(w http.ResponseWriter, r *http.Request) {

}

func (s Service) DeleteMesage(w http.ResponseWriter, r *http.Request) {

}

func (s Service) GetUser(w http.ResponseWriter, r *http.Request) {

}
