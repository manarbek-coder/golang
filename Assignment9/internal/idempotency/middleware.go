package idempotency

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
)

type Store struct {
	DB *sql.DB
}

func NewStore(db *sql.DB) (*Store, error) {
	query := `
	CREATE TABLE IF NOT EXISTS idempotency_keys (
		key TEXT PRIMARY KEY,
		status TEXT NOT NULL,
		status_code INTEGER,
		response_body BLOB
	);`

	_, err := db.Exec(query)
	if err != nil {
		return nil, err
	}

	return &Store{DB: db}, nil
}

func Middleware(store *Store, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Idempotency-Key")
		if key == "" {
			http.Error(w, "Idempotency-Key header required", http.StatusBadRequest)
			return
		}

		result, err := store.DB.Exec(
			"INSERT OR IGNORE INTO idempotency_keys(key, status) VALUES(?, ?)",
			key,
			"processing",
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rows, err := result.RowsAffected()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if rows == 0 {
			var status string
			var statusCode sql.NullInt64
			var responseBody []byte

			err = store.DB.QueryRow(
				"SELECT status, status_code, response_body FROM idempotency_keys WHERE key = ?",
				key,
			).Scan(&status, &statusCode, &responseBody)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if status == "processing" {
				http.Error(w, "Duplicate request in progress", http.StatusConflict)
				return
			}

			if status == "completed" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(int(statusCode.Int64))
				w.Write(responseBody)
				return
			}
		}

		recorder := httptest.NewRecorder()
		next.ServeHTTP(recorder, r)

		body := recorder.Body.Bytes()

		_, err = store.DB.Exec(
			"UPDATE idempotency_keys SET status = ?, status_code = ?, response_body = ? WHERE key = ?",
			"completed",
			recorder.Code,
			body,
			key,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for k, values := range recorder.Header() {
			for _, value := range values {
				w.Header().Add(k, value)
			}
		}

		w.WriteHeader(recorder.Code)
		w.Write(bytes.Clone(body))
	})
}
