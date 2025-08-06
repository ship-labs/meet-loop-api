package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ship-labs/meet-loop-api/internal"
)

type Handler func(w http.ResponseWriter, r *http.Request) Handler

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if next := h(w, r); next != nil {
		next.ServeHTTP(w, r)
	}
}

type Response struct {
	Error   string `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

func OK(w http.ResponseWriter, r *http.Request) Handler {
	return nil
}

func FatalError(v Response, err error) Handler {
	return func(w http.ResponseWriter, r *http.Request) Handler {
		slog.ErrorContext(r.Context(), "fatal error", "origin", "JSON > json.NewEncoder", "data", v, "message", err.Error())
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`Internal server error. Please try again or contact <a href="mailto:dev@shiplabs.dev">Support</a>`))
		return OK
	}
}

func Code(code int, next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) Handler {
		w.WriteHeader(code)
		return next
	}
}

func Text(s string) Handler {
	return func(w http.ResponseWriter, r *http.Request) Handler {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(w, s)
		return OK
	}
}

func CodeText(code int, text string) Handler {
	return Code(code, Text(text))
}

func Error(err error) Handler {
	var code int
	var v internal.ValidationError
	var data any

	switch {
	case errors.As(err, &v):
		code = http.StatusUnprocessableEntity
		data = v.RawErrors()
	case errors.Is(err, internal.ErrExists):
		code = http.StatusConflict
	case errors.Is(err, internal.ErrInvalidRequest):
		code = http.StatusBadRequest
	case errors.Is(err, internal.ErrUnmarshall):
		code = http.StatusBadRequest
	case errors.Is(err, internal.ErrNotExist):
		code = http.StatusNotFound
	case errors.Is(err, internal.ErrGatewayError):
		code = http.StatusBadGateway
	case errors.Is(err, internal.ErrUnauthorized):
		code = http.StatusUnauthorized
	case errors.Is(err, internal.ErrForbidden):
		code = http.StatusForbidden
	default:
		code = http.StatusInternalServerError
	}

	return func(w http.ResponseWriter, r *http.Request) Handler {
		if code == http.StatusInternalServerError {
			slog.Log(r.Context(), slog.LevelError, "internal", "url", r.URL.Path, "error", err)
			err = internal.ErrInternal
		}

		if code == http.StatusBadRequest {
			slog.Log(r.Context(), slog.LevelError, "bad request", "url", r.URL.Path, "error", err)
			err = internal.ErrInvalidRequest
		}

		return Code(code, JSON(Response{
			Error:   err.Error(),
			Message: http.StatusText(code),
			Data:    data,
		}))
	}
}

func JSON(v Response) Handler {
	return func(w http.ResponseWriter, r *http.Request) Handler {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(v); err != nil {
			return FatalError(v, err)
		}
		return OK
	}
}
