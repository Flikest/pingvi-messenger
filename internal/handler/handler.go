package handler

import (
	"github.com/Flikest/PingviMessenger/internal/services"
	"github.com/gorilla/mux"
)

type Handler struct {
	Service *services.Service
}

func InitHandler(s *services.Service) *Handler {
	return &Handler{Service: s}
}

func (h Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/mesenger", h.Service.Correspondence)
	return r
}
