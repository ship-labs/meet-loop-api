// Package internal provides shared validation utilities for the IELTS Agent application.
package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/Oudwins/zog"
)

type ValidationError zog.ZogIssueMap

func (v ValidationError) Error() string {
	var errors []string
	for _, value := range v {
		for _, issue := range value {
			errors = append(errors, fmt.Sprintf("%s: %s", issue.Path, issue.Message))
		}
	}

	return strings.Join(errors, ", ")
}

func ValidateJSON[T any](r io.Reader, validator func(p *T) error) (T, error) {
	var data T

	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&data)
	if err != nil {
		return data, fmt.Errorf("%w: unmarshalling json: %w", ErrUnmarshall, err)
	}

	// Validate and parse data
	err = validator(&data)
	if e, ok := err.(ValidationError); ok {
		return data, e
	}

	if err != nil {
		return data, fmt.Errorf("%w: %w: ", ErrInvalidRequest, err)
	}

	return data, nil
}
