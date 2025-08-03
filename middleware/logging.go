package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"time"
)

var (
	assetExts = []string{".css", ".js", ".svg", ".png", ".jpg"}
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(statusCode int) {
	rr.statusCode = statusCode
	rr.ResponseWriter.WriteHeader(statusCode)
}

func (rr *responseRecorder) Write(p []byte) (int, error) {
	if rr.statusCode == 0 {
		rr.statusCode = http.StatusOK
	}
	return rr.ResponseWriter.Write(p)
}

// This method is special and can be detected later by handlers to allow them to use optional
// interfaces without using type assertion.
func (rr *responseRecorder) Unwrap() http.ResponseWriter {
	return rr.ResponseWriter
}

func LoggingMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rr := &responseRecorder{ResponseWriter: w}
		next.ServeHTTP(rr, r)

		if r.Method == http.MethodGet {
			for _, ext := range assetExts {
				if strings.HasSuffix(strings.ToLower(r.URL.Path), ext) {
					return
				}
			}
		}

		slog.Log(r.Context(), slog.LevelInfo, "request",
			"url", r.URL.String(),
			"method", r.Method,
			"took (ms)", time.Since(start).Milliseconds(),
			"statusCode", rr.statusCode,
			"userAgent", r.UserAgent(),
			"remoteAddr", r.RemoteAddr,
			"X-Forwarded-For", r.Header.Get("X-Forwarded-For"),
		)
	}
}
