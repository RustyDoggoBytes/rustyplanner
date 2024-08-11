package main

import (
	"crypto/subtle"
	"net/http"
	"os"
	"time"
)

func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		value = fallback
	}
	return value
}

func getMondayOfCurrentWeek(now time.Time) time.Time {
	weekday := now.Weekday()
	var monday time.Time
	if weekday == time.Sunday {
		// If today is Sunday, we need to go back 6 days
		monday = now.AddDate(0, 0, -6)
	} else {
		// Otherwise, we go back (weekday - 1) days
		monday = now.AddDate(0, 0, -int(weekday)+1)
	}
	return monday.Truncate(24 * time.Hour)
}

func basicAuthMiddleware(next http.Handler, username, password string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
