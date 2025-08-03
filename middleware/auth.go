// Package middleware provides HTTP middleware for authentication and related utilities.
package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ship-labs/meet-loop-api/config"
	"github.com/ship-labs/meet-loop-api/internal"
)

// Define an unexported type for context keys to avoid collisions
type contextKey struct{}

// Define specific keys as instances of the empty struct
var jwtClaimsKey = contextKey{}

// UserMetadata represents the structure of the user_metadata in the JWT
type UserMetadata struct {
	BandListening string `json:"band_listening"`
	BandReading   string `json:"band_reading"`
	BandSpeaking  string `json:"band_speaking"`
	BandWriting   string `json:"band_writing"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	EmailVerified bool   `json:"email_verified"`
	PhoneVerified bool   `json:"phone_verified"`
	Sub           string `json:"sub"`
}

// AppMetadata represents the structure of the app_metadata in the JWT
type AppMetadata struct {
	Provider  string   `json:"provider"`
	Providers []string `json:"providers"`
}

// AMR represents an authentication method record in the JWT
type AMR struct {
	Method    string `json:"method"`
	Timestamp int64  `json:"timestamp"`
}

// JWTClaims represents the structure of the JWT payload
type JWTClaims struct {
	Iss          string       `json:"iss"`
	Sub          string       `json:"sub"`
	Aud          string       `json:"aud"`
	Exp          int64        `json:"exp"`
	Iat          int64        `json:"iat"`
	Email        string       `json:"email"`
	Phone        string       `json:"phone"`
	AppMetadata  AppMetadata  `json:"app_metadata"`
	UserMetadata UserMetadata `json:"user_metadata"`
	Role         string       `json:"role"`
	AAL          string       `json:"aal"`
	AMR          []AMR        `json:"amr"`
	SessionID    string       `json:"session_id"`
	IsAnonymous  bool         `json:"is_anonymous"`
	jwt.RegisteredClaims
}

func Auth(next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) Handler {
		cfg, err := config.LoadConfig()
		if err != nil {
			return Error(err)
		}

		auth := r.Header.Get("Authorization")
		prefix, tokenString, ok := strings.Cut(auth, " ")
		if !ok {
			return Error(internal.ErrUnauthorized)
		}

		if !strings.EqualFold(prefix, "Bearer") {
			return Error(internal.ErrInvalidRequest)
		}

		//   Parse token with custom claims struct
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.JwtSecret), nil
		})

		if err != nil {
			return Error(fmt.Errorf("parsing token: %w", err))
		}

		if !token.Valid {
			return Error(internal.ErrUnauthorized)
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			return Error(fmt.Errorf("invalid token claims"))
		}

		// Update the request context with the claims
		ctx := context.WithValue(r.Context(), jwtClaimsKey, claims)

		// Call the next handler with the updated request
		return next(w, r.WithContext(ctx))
	}
}

// GetClaims returns the full JWT claims object from the context
func GetClaims(ctx context.Context) (*JWTClaims, error) {
	claims, ok := ctx.Value(jwtClaimsKey).(*JWTClaims)
	if !ok {
		return nil, fmt.Errorf("no JWT claims found in context")
	}

	return claims, nil
}

// GetUserID extracts the user ID from the JWT claims in the context
func GetUserID(ctx context.Context) (pgtype.UUID, error) {
	claims, err := GetClaims(ctx)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("getting userID: %w", err)
	}

	userID, err := uuid.Parse(claims.Sub)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("parsing userID %s: %w", claims.Sub, err)
	}

	return pgtype.UUID{Bytes: userID, Valid: true}, nil
}

// GetUserEmail extracts the user email from the JWT claims in the context
func GetUserEmail(ctx context.Context) (string, error) {
	claims, err := GetClaims(ctx)
	if err != nil {
		return "", fmt.Errorf("getting email: %w", err)
	}

	return claims.Email, nil
}

// GetUserMetadata extracts the user metadata from the JWT claims in the context
func GetUserMetadata(ctx context.Context) (UserMetadata, error) {
	claims, err := GetClaims(ctx)
	if err != nil {
		return UserMetadata{}, fmt.Errorf("getting user metadata: %w", err)
	}

	return claims.UserMetadata, nil
}
