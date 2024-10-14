package config

import (
	"os"
)

func newEnv() preConfig {
	pre := newPreConfig()

	salt, exists := os.LookupEnv("SALT")
	if exists {
		pre.Salt = &salt
	}

	jwt, exists := os.LookupEnv("JWT_KEY")
	if exists {
		pre.JwtKey = &jwt
	}

	return pre
}
