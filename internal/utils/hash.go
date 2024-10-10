package utils

import (
	"crypto/sha512"
	"encoding/hex"
	"slices"
)

func GetHash(s string, salt string) string {
	bs := slices.Concat([]byte(s), []byte(salt))
	hasher := sha512.New()
	hasher.Write(bs)
	return hex.EncodeToString(hasher.Sum(nil))
}
