// Package internal contains shared constants and utilities for the IELTS Agent application.
package internal

type Achievement struct {
	Title    string `json:"title"`
	Icon     string `json:"icon"`
	Unlocked bool   `json:"unlocked"`
}

const (
	UniqueViolationCode = "23505"
	InvestorsLimit      = 150
	ProgramRoot         = "cmd"
)

var Achievements = []Achievement{
	{Title: "First Test Completed", Icon: "🎯", Unlocked: true},
	{Title: "High Scorer", Icon: "⭐", Unlocked: false},
	{Title: "Practice Master", Icon: "🏆", Unlocked: false},
}
