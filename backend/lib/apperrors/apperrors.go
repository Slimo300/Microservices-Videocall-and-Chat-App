package apperrors

import (
	"errors"
	"log"
	"net/http"
)

type Error struct {
	statusCode int
	message    string
}

func (e *Error) Error() string { return e.message }

func (e *Error) StatusCode() int { return e.statusCode }

func Status(err error) int {
	var e *Error
	if errors.As(err, &e) {
		return e.StatusCode()
	}
	return http.StatusInternalServerError
}

func NewInternal(err error) *Error {
	log.Println(err.Error())
	return &Error{
		statusCode: http.StatusInternalServerError,
		message:    "internal server error",
	}
}

func NewAuthorization(message string) *Error {
	return &Error{
		statusCode: http.StatusUnauthorized,
		message:    message,
	}
}

func NewBadRequest(message string) *Error {
	return &Error{
		statusCode: http.StatusBadRequest,
		message:    message,
	}
}

func NewForbidden(message string) *Error {
	return &Error{
		statusCode: http.StatusForbidden,
		message:    message,
	}
}

func NewConflict(message string) *Error {
	return &Error{
		statusCode: http.StatusConflict,
		message:    message,
	}
}

func NewNotFound(message string) *Error {
	return &Error{
		statusCode: http.StatusNotFound,
		message:    message,
	}
}

func NewPayloadTooLarge(message string) *Error {
	return &Error{
		statusCode: http.StatusRequestEntityTooLarge,
		message:    message,
	}
}

func NewUnsupportedMediaType(message string) *Error {
	return &Error{
		statusCode: http.StatusUnsupportedMediaType,
		message:    message,
	}
}
