package auth

import (
	"fmt"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		fmt.Println("Error when hashing:", err)
		return "", err
	}

	return hash, nil
}

func CheckPasswordHash(password string, hash string) (bool, error) {
	doesMatch, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		fmt.Println("Error while comparing hash:", err)
		return false, err
	}

	return doesMatch, err
}
