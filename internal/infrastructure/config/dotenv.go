package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

func newDotenv() preConfig {
	err := godotenv.Load(".env.local")
	if err != nil {
		_ = godotenv.Load(".env")
	}

	return preConfig{
		Salt:   getEnv("SALT", ""),
		JwtKey: getEnv("JWT_KEY", ""),
	}
}

func getEnvInt(key string, def int) *int {
	v, e := strconv.Atoi(*getEnv(key, strconv.Itoa(def)))
	if e != nil {
		return &def
	} else {
		return &v
	}
}

func getEnv(key string, def string) *string {
	if value, exists := os.LookupEnv(key); exists {
		return &value
	}
	return &def
}
