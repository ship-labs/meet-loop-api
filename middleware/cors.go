package middleware

import (
	"net/http"

	"github.com/ship-labs/meet-loop-api/config"
)

func CorsMiddleware(next http.Handler) http.Handler {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers conditionally
		w.Header().Set("Access-Control-Allow-Origin", cfg.FrontendURL)
		w.Header().Set("Vary", "Origin")

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
