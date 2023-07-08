package goerrors

import "net/http"

type UnauthorizedError struct {
	message    string
	StatusCode int
}

func (ue UnauthorizedError) Error() string {
	return ue.message
}

func (ue UnauthorizedError) GetStatusCode() int {
	return ue.StatusCode
}

func NewUnauthorizedError(message string) UnauthorizedError {
	return UnauthorizedError{
		message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}
