// Package internal provides shared validation utilities for the IELTS Agent application.
package internal

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/Oudwins/zog"
	"github.com/Oudwins/zog/zconst"
)

type ValidationError zog.ZogIssueMap

func (v ValidationError) Error() string {
	return "Incorrect or missing form data."
}

func (v ValidationError) RawErrors() []string {
	var errors []string
	for key, values := range zog.Issues.SanitizeMap(v) {
		if key != zconst.ISSUE_KEY_FIRST {
			errors = append(errors, values...)
		}
	}
	return errors
}

func Validate[T any](schema *zog.StructSchema, body io.Reader) (result T, err error) {
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&result); err != nil {
		return result, fmt.Errorf("%w:decoding json: %w", ErrInvalidRequest, err)
	}

	if err := schema.Validate(&result); len(err) != 0 {
		return result, ValidationError(err)
	}

	return result, nil
}
