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

	return pre
}
