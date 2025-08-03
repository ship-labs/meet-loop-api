package middleware

import (
	"encoding/json"
	"io"
	"log/slog"

	"github.com/Oudwins/zog"
	"github.com/ship-labs/meet-loop-api/internal"
)

func Validate[T any](schema *zog.StructSchema, body io.Reader) (result T, err error) {
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		slog.Error("validating json", "error", err)
		return result, internal.ErrInvalidRequest
	}

	if err := schema.Validate(&result); len(err) != 0 {
		return result, internal.ValidationError(err)
	}

	return result, nil
}
