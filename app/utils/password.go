package utils

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	PW_SALT_BYTES = 32
	PW_HASH_BYTES = 64
)

func HashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}

func CheckPassword(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
