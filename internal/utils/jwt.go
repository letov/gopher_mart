package utils

import (
	"context"
	"errors"
	"github.com/go-chi/jwtauth"
	"time"
)

var (
	ErrInvalidUserId = errors.New("invalid user id")
)

func GetJwt(secret string, userID int64) (string, error) {
	createAt := time.Now().Unix()
	expireAt := time.Now().Add(time.Minute * 10)
	tokenAuth := jwtauth.New("HS256", []byte(secret), nil)

	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
		"c_at":    createAt,
		"e_at":    expireAt,
		"user_id": userID,
	})
	return tokenString, err
}

func GetUserIdFromToken(ctx context.Context) (int64, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return 0, err
	}

	if _, ok := claims["user_id"]; !ok {
		return 0, ErrInvalidUserId
	}

	userId, ok := claims["user_id"].(float64)
	if !ok {
		return 0, ErrInvalidUserId
	}

	return int64(userId), nil
}
