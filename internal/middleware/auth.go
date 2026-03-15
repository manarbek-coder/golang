package middleware

import (
<<<<<<< HEAD
	"net/http"
=======
	"log"
	"net/http"
	"time"
>>>>>>> d902c865a99a662e7408aa1cecd6b5830edb5ea5
)

const validAPIKey = "secret12345"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
<<<<<<< HEAD
=======
		start := time.Now()
		log.Printf("%s %s", r.Method, r.URL.Path)

>>>>>>> d902c865a99a662e7408aa1cecd6b5830edb5ea5
		apiKey := r.Header.Get("X-API-KEY")
		if apiKey != validAPIKey {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"unauthorized"}`))
			return
		}

		next.ServeHTTP(w, r)
<<<<<<< HEAD
=======
		log.Printf("%s %s completed in %v", r.Method, r.URL.Path, time.Since(start))
>>>>>>> d902c865a99a662e7408aa1cecd6b5830edb5ea5
	})
}
