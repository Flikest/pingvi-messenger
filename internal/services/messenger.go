package services

import (
	"net/http"

	"github.com/Flikest/PingviMessenger/internal/storage"

	"github.com/gorilla/websocket"
)

type Service struct {
	Storage *storage.Storage
}

var upgreader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func InitService(s *storage.Storage) *Service {
	return &Service{Storage: s}
}

func (s Service) Correspondence(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
	}
}

func (s Service) UpdateMessage(w http.ResponseWriter, r *http.Request) {

}

func (s Service) DeleteMesage(w http.ResponseWriter, r *http.Request) {

}

func (s Service) GetUser(w http.ResponseWriter, r *http.Request) {

}
