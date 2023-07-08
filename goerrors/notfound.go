package goerrors

import "net/http"

type NotFoundError struct {
	message    string
	StatusCode int
}

func (nfe NotFoundError) Error() string {
	return nfe.message
}

func (nfe NotFoundError) GetStatusCode() int {
	return nfe.StatusCode
}

func NewNotFoundError(message string) NotFoundError {
	return NotFoundError{
		message:    message,
		StatusCode: http.StatusNotFound,
	}
}
