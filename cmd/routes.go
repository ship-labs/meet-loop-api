package main

import (
	"net/http"

	"github.com/ship-labs/meet-loop-api/internal"
	"github.com/ship-labs/meet-loop-api/internal/pkg/sqlc"
	"github.com/ship-labs/meet-loop-api/members"
	"github.com/ship-labs/meet-loop-api/middleware"
)

func defineRoutes(store *sqlc.Store) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /{$}", middleware.Auth(func(w http.ResponseWriter, r *http.Request) middleware.Handler {
		return middleware.JSON(middleware.Response{Message: http.StatusText(http.StatusOK)})
	}))

	mux.Handle(internal.APIVersion, middleware.JSON(middleware.Response{
		Message: http.StatusText(http.StatusOK),
		Data:    "Welcome to MeetLoop API v1",
	}))

	mux.Handle(internal.Group, middleware.Auth(members.CreateGroup(store)))
	mux.Handle(internal.Profile, middleware.Auth(members.GetUserProfile(store)))

	return mux
}
