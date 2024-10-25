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

	ampq, exists := os.LookupEnv("AMPQ")
	if exists {
		pre.Ampq = &ampq
	}

	accrual, exists := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS")
	if exists {
		pre.AccrualUrl = &accrual
	}

	connection, exists := os.LookupEnv("DATABASE_URI")
	if exists {
		pre.DBConnection = &connection
	}

	addr, exists := os.LookupEnv("RUN_ADDRESS")
	if exists {
		pre.Addr = &addr
	}

	return pre
}
