package utils

import (
	"github.com/go-chi/jwtauth"
	"gopher_mart/internal/application/dto/result"
	"time"
)

func GetJwt(secret string, login result.Login) (string, error) {
	createAt := time.Now().Unix()
	expireAt := time.Now().Add(time.Minute * 10)
	tokenAuth := jwtauth.New("HS256", []byte(secret), nil)

	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
		"c_at":    createAt,
		"e_at":    expireAt,
		"user_id": login.UserID,
	})
	return tokenString, err
}
