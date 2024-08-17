package middlewares

import (
	"log/slog"
	"net/http"
)

type CustomResponseWriter struct {
	http.ResponseWriter
	status int
}

func NewCustomResponseWriter(w http.ResponseWriter) *CustomResponseWriter {
	return &CustomResponseWriter{w, http.StatusOK}
}

func (crw *CustomResponseWriter) WriteHeader(code int) {
	crw.status = code
	crw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a custom ResponseWriter to capture the status code
		crw := NewCustomResponseWriter(w)

		// Start timer
		//start := time.Now()

		// Call the next handler
		next.ServeHTTP(crw, r)

		// Calculate duration
		//duration := time.Since(start)

		// Log the request
		slog.Info("HTTP",
			"method", r.Method,
			"path", r.URL.Path,
			"status", crw.status,
			//"duration", duration,
			//"ip", r.RemoteAddr,
			//"user_agent", r.UserAgent(),
		)
	})
}
