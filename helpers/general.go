package helpers

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func CheckPassword(password, hash string) error {
	fmt.Println("check password", password, hash)
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}