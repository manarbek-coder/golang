package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("[%s] %s %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("[%s] %s %s completed in %v",
			time.Now().Format(time.RFC3339),
			r.Method,
			r.URL.Path,
			time.Since(start))
	})
}
