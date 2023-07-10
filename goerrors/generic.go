package goerrors

// GenericHttpError defines an interface to handle custom http errors through out the code
type GenericHttpError interface {
	Error() string
	GetStatusCode() int
}
