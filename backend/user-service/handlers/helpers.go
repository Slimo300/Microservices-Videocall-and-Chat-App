package handlers

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

var ErrNoDatabase = errors.New("No database connection")
var InvalidToken = errors.New("Invalid token")

func isPasswordValid(pass string) bool {
	return len(pass) >= 6
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func hashPassword(s string) (string, error) {
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

func checkPassword(existingHash, incomingPass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(existingHash), []byte(incomingPass)) == nil
}
