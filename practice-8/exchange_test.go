package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetRate(t *testing.T) {
	tests := []struct {
		name      string
		handler   http.HandlerFunc
		expectErr bool
	}{
		{
			name: "Successful scenario",
			handler: func(w http.ResponseWriter, r *http.Request) {
				resp := RateResponse{Base: "USD", Target: "EUR", Rate: 0.85}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(resp)
			},
			expectErr: false,
		},
		{
			name: "API Business Error - 404 with error message",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{"error": "invalid currency pair"})
			},
			expectErr: true,
		},
		{
			name: "Malformed JSON",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal Server Error"))
			},
			expectErr: true,
		},
		{
			name: "Slow Response Timeout",
			handler: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(6 * time.Second)
				w.WriteHeader(http.StatusOK)
			},
			expectErr: true,
		},
		{
			name: "Server Panic 500",
			handler: func(w http.ResponseWriter, r *http.Request) {
				panic("server panic")
			},
			expectErr: true,
		},
		{
			name: "Empty Body",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			service := NewExchangeService(server.URL)
			_, err := service.GetRate("USD", "EUR")

			if tt.expectErr && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
		})
	}
}
