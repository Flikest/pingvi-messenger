package services

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/Flikest/PingviMessenger/internal/storage"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

	"github.com/gorilla/websocket"
)

type Service struct {
	Storage *storage.Storage
}

type tokenClaims struct {
	jwt.StandardClaims
	ID uuid.UUID `json:"id"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func InitService(s *storage.Storage) *Service {
	return &Service{Storage: s}
}

func (s Service) Correspondence(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Info("пиздык системе, переделывай!!!: ", err)
	}
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			slog.Info("и тут пиде, а что ты хотел, стив джобс хуев:", err)
			return
		}
		// print out that message for clarity
		slog.Info(string(p))
		if err := conn.WriteMessage(messageType, p); err != nil {
			slog.Info("ну давай блять", err)
			return
		}
	}
}

func ParseTocken(accessToken string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return uuid.Nil, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.ID, nil
}
