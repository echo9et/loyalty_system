package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sha256hash(value string) string {
	hash := sha256.New()
	hash.Write([]byte(value))
	return hex.EncodeToString(hash.Sum(nil))
}
