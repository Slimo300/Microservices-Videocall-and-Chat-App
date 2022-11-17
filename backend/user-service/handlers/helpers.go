package handlers

import (
	"errors"
	"regexp"
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

var validImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
}

// IsAllowedImageType determines if image is among types defined
// in map of allowed images
func isAllowedImageType(mimeType string) bool {
	_, exists := validImageTypes[mimeType]

	return exists
}
