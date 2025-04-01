package jwt

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/golang-jwt/jwt"
)

func JwtPayloadFromRequest(tokenString string) (string, error) {

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return "", err
	}

	result, err := json.Marshal(claims["sub"])
	if err != nil {
		slog.Info("Error: %s", err)
		return "", err
	}
	return string(result), nil
}
