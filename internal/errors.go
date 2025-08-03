package internal

import (
	"errors"
	"net/http"
)

var (
	ErrExists         = errors.New("already exists")
	ErrNotExist       = errors.New("does not exist")
	ErrInvalidRequest = errors.New("invalid request")
	ErrGatewayError   = errors.New("error")
	ErrInternal       = errors.New(
		"internal error: please try again later or contact support",
	)
	ErrUnmarshall   = errors.New("unmarshalling json error")
	ErrUnauthorized = errors.New(http.StatusText(http.StatusUnauthorized))
	ErrForbidden    = errors.New(http.StatusText(http.StatusForbidden))
)
