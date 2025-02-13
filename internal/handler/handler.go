package handler

import "github.com/Flikest/PingviMessenger/internal/services"

type Handler struct {
	Service *services.Service
}

func InitHandler(s *services.Service) *Handler {
	return &Handler{Service: s}
}
