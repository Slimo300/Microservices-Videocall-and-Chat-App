package apperrors

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Error is a wrapper around error message that holds
// error type (its http status)
type Error struct {
	statusCode int
	message    string
}

// Error fulfills error interface
func (e *Error) Error() string { return e.message }

// Status returns http status code associated with error
func (e *Error) StatusCode() int { return e.statusCode }

// Status returns corresponding http status if err is of type Error
// if not it returns status 500
func Status(err error) int {
	var e *Error
	if errors.As(err, &e) {
		return e.StatusCode()
	}
	return http.StatusInternalServerError
}

// NewAuthorization creates Authorization Error with given message
func NewAuthorization(message string) *Error {
	return &Error{
		statusCode: http.StatusUnauthorized,
		message:    message,
	}
}

// NewBadRequest creates BadRequest Error with given message
func NewBadRequest(message string) *Error {
	return &Error{
		statusCode: http.StatusBadRequest,
		message:    message,
	}
}

// NewForbidden creates Forbidden Error with given message
func NewForbidden(message string) *Error {
	return &Error{
		statusCode: http.StatusForbidden,
		message:    message,
	}
}

// NewConflict constructs Conflict Error with given parameters
func NewConflict(message string) *Error {
	return &Error{
		statusCode: http.StatusConflict,
		message:    message,
	}
}

// NewNotFound constructs new NotFound Error with given parameters
func NewNotFound(message string) *Error {
	return &Error{
		statusCode: http.StatusNotFound,
		message:    message,
	}
}

// NewPayloadTooLarge constructs new PayloadTooLarge error with given parameters
func NewPayloadTooLarge(message string) *Error {
	return &Error{
		statusCode: http.StatusRequestEntityTooLarge,
		message:    message,
	}
}

// NewUnsupportedMediaType creates new UnsupportedMediaType Error with given message
func NewUnsupportedMediaType(message string) *Error {
	return &Error{
		statusCode: http.StatusUnsupportedMediaType,
		message:    message,
	}
}

type GinErrorHandler struct {
	logger *log.Logger
}

func NewGinErrorHandler() *GinErrorHandler {
	return &GinErrorHandler{}
}

// WithLogger appends logger to our handler
func (eh *GinErrorHandler) WithLogger(logger *log.Logger) *GinErrorHandler {
	eh.logger = logger
	return eh
}

func (eh *GinErrorHandler) HandleError(c *gin.Context, err error) {
	if eh.logger != nil {
		eh.logger.Println(err.Error())
	}

	c.JSON(Status(err), err.Error())
	return
}
