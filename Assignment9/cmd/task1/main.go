package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"time"

	"Assignment9/internal/retry"
)

func main() {
	var counter int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		current := atomic.AddInt32(&counter, 1)

		if current <= 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"error":"service unavailable"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := retry.PaymentClient{
		Client:     &http.Client{},
		MaxRetries: 5,
	}

	body, err := client.ExecutePayment(ctx, server.URL)
	if err != nil {
		fmt.Println("Payment failed:", err)
		return
	}

	fmt.Println("Final response:", string(body))
}
