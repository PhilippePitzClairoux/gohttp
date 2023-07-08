package goerrors

import "net/http"

type BadRequestError struct {
	message    string
	StatusCode int
}

func (bre BadRequestError) Error() string {
	return bre.message
}

func (bre BadRequestError) GetStatusCode() int {
	return bre.StatusCode
}

func NewBadRequestError(message string) BadRequestError {
	return BadRequestError{
		message:    message,
		StatusCode: http.StatusBadRequest,
	}
}
