package services

import "github.com/Flikest/PingviMessenger/internal/storage"

type Service struct {
	Storage *storage.Storage
}

func NewServices(s *storage.Storage) *Service {
	return &Service{Storage: s}
}
