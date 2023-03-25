package database

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(s string) (string, error) {
	if s == "" {
		return "", errors.New("Reference provided for hashing password is nil")
	}
	sBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(sBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	s = string(hashedBytes)
	return s, nil
}

func CheckPassword(existingHash, incomingPass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(existingHash), []byte(incomingPass)) == nil
}
