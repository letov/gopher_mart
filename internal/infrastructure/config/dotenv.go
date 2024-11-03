package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

func newDotenv() preConfig {
	var err error
	if os.Getenv("IS_TEST_ENV") == "true" {
		err = godotenv.Load("../../.env.test")
	} else {
		err = godotenv.Load(".env.local")
		if err != nil {
			err = godotenv.Load(".env")
		}
	}

	return preConfig{
		Salt:         getEnv("SALT", ""),
		JwtKey:       getEnv("JWT_KEY", ""),
		Ampq:         getEnv("AMPQ", ""),
		AccrualUrl:   getEnv("ACCRUAL_SYSTEM_ADDRESS", ""),
		DBConnection: getEnv("DATABASE_URI", ""),
		Addr:         getEnv("RUN_ADDRESS", ""),
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
