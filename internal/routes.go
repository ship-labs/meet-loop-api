// Package internal contains shared constants and utilities for the IELTS Agent application.
package internal

import (
	"fmt"
	"net/http"
)

var (
	APIVersion                = "/api/v1"
	GetDashboardData          = createRoute(http.MethodGet, "get-dashboard-data")
	SendMessage               = createRoute(http.MethodPost, "chat")
	GetMessages               = createRoute(http.MethodGet, "chat")
	GetPractiseTests          = createRoute(http.MethodGet, "practice-tests")
	GetTestQuestions          = createRoute(http.MethodGet, "questions/{testID}")
	SubmitTest                = createRoute(http.MethodPost, "practise-tests/{testID}")
	DeleteShortlistedInvestor = createRoute(http.MethodDelete, "shortlist/{investorID}")
)

func createRoute(method, path string) string {
	return fmt.Sprintf("%s %s/%s", method, APIVersion, path)
}
