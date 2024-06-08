package service

import "net/http"

type Error struct {
	Err     error
	Status  int
	Message string
}

func internalServerError(err error, msg string) *Error {
	return &Error{
		Err:     err,
		Status:  http.StatusInternalServerError,
		Message: msg,
	}
}

func badRequest(err error, msg string) *Error {
	return &Error{
		Err:     err,
		Status:  http.StatusBadRequest,
		Message: msg,
	}
}

func unauthorized(err error, msg string) *Error {
	return &Error{
		Err:     err,
		Status:  http.StatusUnauthorized,
		Message: msg,
	}
}

func otherError(err error, msg string, status int) *Error {
	return &Error{
		Err:     err,
		Message: msg,
		Status:  status,
	}
}
