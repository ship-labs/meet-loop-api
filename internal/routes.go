// Package internal contains shared constants and utilities for the IELTS Agent application.
package internal

import (
	"fmt"
	"net/http"
	"strings"
)

var (
	APIVersion = "/api/v1"
	Group      = createRoute(http.MethodPost, "group")
	Profile    = createRoute(http.MethodGet, "/profile")
)

func createRoute(method, path string) string {
	if strings.HasPrefix(path, "/") {
		_, path, _ = strings.Cut(path, "/")
	}
	return fmt.Sprintf("%s %s/%s", method, APIVersion, path)
}
