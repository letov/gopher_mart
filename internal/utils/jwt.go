package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GetJwt(secret string) (string, error) {
	createAt := time.Now().Unix()
	expireAt := time.Now().Add(time.Minute * 10)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"c_at": createAt,
		"e_at": expireAt,
	})
	return token.SignedString([]byte(secret))
}
