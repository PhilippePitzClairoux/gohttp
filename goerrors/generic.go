package goerrors

type GenericHttpError interface {
	Error() string
	GetStatusCode() int
}
