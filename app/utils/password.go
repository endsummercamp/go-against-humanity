package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

func HashPassword(password string) string {
	hasher := sha256.New()
	io.WriteString(hasher, password)
	return hex.EncodeToString(hasher.Sum(nil))
}
