package util

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang-clean-architecture/internal/model"
	"time"
)

type TokenUtil struct {
	SecretKey string
}

func NewTokenUtil(secretKey string) *TokenUtil {
	return &TokenUtil{
		SecretKey: secretKey,
	}
}

func (t TokenUtil) CreateToken(auth *model.Auth) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":     auth.ID,
		"expire": time.Now().Add(time.Hour * 24 * 30).UnixMilli(),
	})

	jwtToken, err := token.SignedString([]byte(t.SecretKey))
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func (t TokenUtil) ParseToken(jwtToken string) (*model.Auth, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.SecretKey), nil
	})
	if err != nil {
		return nil, fiber.ErrUnauthorized
	}

	claims := token.Claims.(jwt.MapClaims)

	expire := claims["expire"].(float64)
	if int64(expire) < time.Now().UnixMilli() {
		return nil, fiber.ErrUnauthorized
	}

	id := claims["id"].(string)
	auth := &model.Auth{
		ID: id,
	}
	return auth, nil
}
