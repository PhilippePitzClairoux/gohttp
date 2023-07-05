package goerrors

import "net/http"

type UnauthorizedError struct {
	message    string
	StatusCode int
}

func (ue UnauthorizedError) Error() string {
	return ue.message
}

func NewUnauthorizedError(message string) UnauthorizedError {
	return UnauthorizedError{
		message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}
