package handler

import (
	"net/http"

	"github.com/Flikest/PingviMessenger/internal/services"
)

type Handler struct {
	Service *services.Service
}

func InitHandler(s *services.Service) *Handler {
	return &Handler{Service: s}
}

func (h Handler) NewRouter() {
	http.HandleFunc("/mesenger", h.Service.Correspondence)
}
