package service

import (
	"errors"
	"github.com/Nerzal/gocloak/v13"
	"net/http"
)

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

func keyCloakError(err error) *Error {
	var apiErr *gocloak.APIError
	if errors.As(err, &apiErr) {
		return &Error{
			Err:     apiErr,
			Message: apiErr.Message,
			Status:  apiErr.Code,
		}
	}
	return internalServerError(err, "unknown")
}
