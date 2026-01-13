package errors

import "net/http"

type Error interface {
	error
	StatusCode() int
}

type serviceError struct {
	msg  string
	code int
}

func (e serviceError) Error() string { //since we are using error interface, we need to implement Error() method
	return e.msg
}

func (e serviceError) StatusCode() int {
	return e.code
}

func NewError(msg string, code int) Error {
	return serviceError{msg: msg, code: code}
}

var (
	ErrInvalidPayload     = NewError("invalid payload", http.StatusBadRequest)
	ErrInvalidCredentials = NewError("invalid email or password", http.StatusBadRequest)
	ErrTokenGeneration    = NewError("failed to generate token", http.StatusInternalServerError)
	ErrTokenGeneration = NewError("failed to generate token", http.StatusInternalServerError
)
