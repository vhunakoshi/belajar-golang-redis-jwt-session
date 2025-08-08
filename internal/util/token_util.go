package util

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang-clean-architecture/internal/model"
	"time"
)

type TokenUtil struct {
	SecretKey string
	Redis     *redis.Client
}

func NewTokenUtil(secretKey string, redistClient *redis.Client) *TokenUtil {
	return &TokenUtil{
		SecretKey: secretKey,
		Redis:     redistClient,
	}
}

func (t TokenUtil) CreateToken(ctx context.Context, auth *model.Auth) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":     auth.ID,
		"expire": time.Now().Add(time.Hour * 24 * 30).UnixMilli(),
	})

	jwtToken, err := token.SignedString([]byte(t.SecretKey))
	if err != nil {
		return "", err
	}

	_, err = t.Redis.SetEx(ctx, jwtToken, auth.ID, time.Hour*25*30).Result()
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func (t TokenUtil) ParseToken(ctx context.Context, tokenString string) (*model.Auth, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.SecretKey), nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil {
		return nil, fiber.ErrUnauthorized
	}

	claims := token.Claims.(jwt.MapClaims)
	expire := int64(claims["expire"].(float64))
	if expire < time.Now().UnixMilli() {
		return nil, fiber.ErrUnauthorized
	}

	result, err := t.Redis.Exists(ctx, tokenString).Result()
	if err != nil {
		return nil, err
	}

	if result == 0 {
		return nil, fiber.ErrUnauthorized
	}
	
	id := claims["id"].(string)
	auth := &model.Auth{
		ID: id,
	}

	return auth, nil
}
