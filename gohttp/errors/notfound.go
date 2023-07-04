package gohttperrors

import "net/http"

type NotFoundError struct {
	message    string
	StatusCode int
}

func (nfe NotFoundError) Error() string {
	return nfe.message
}

func NewNotFoundError(message string) NotFoundError {
	return NotFoundError{
		message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}
