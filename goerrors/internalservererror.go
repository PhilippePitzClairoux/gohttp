package goerrors

import "net/http"

type InternalServerError struct {
	message    string
	StatusCode int
}

func (ise InternalServerError) Error() string {
	return ise.message
}

func (ise InternalServerError) GetStatusCode() int {
	return ise.StatusCode
}

func NewInternalServerError(message string) InternalServerError {
	return InternalServerError{
		message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}
