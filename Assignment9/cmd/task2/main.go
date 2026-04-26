package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"

	"Assignment9/internal/idempotency"
)

func main() {
	db, err := sql.Open("sqlite", "file:idempotency.db?cache=shared")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	store, err := idempotency.NewStore(db)
	if err != nil {
		panic(err)
	}

	var businessCalls int32

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&businessCalls, 1)

		fmt.Println("Processing started")
		time.Sleep(2 * time.Second)

		transactionID := "uuid-" + uuid.NewString()
		body := fmt.Sprintf(`{"status":"paid","amount":1000,"transaction_id":"%s"}`, transactionID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))

		fmt.Println("Processing completed")
	})

	server := httptest.NewServer(idempotency.Middleware(store, handler))
	defer server.Close()

	key := "same-payment-key-1000"
	client := &http.Client{}
	var wg sync.WaitGroup

	for i := 1; i <= 10; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			req, err := http.NewRequest(http.MethodPost, server.URL, bytes.NewBuffer(nil))
			if err != nil {
				fmt.Println("Request", id, "error:", err)
				return
			}

			req.Header.Set("Idempotency-Key", key)

			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Request", id, "error:", err)
				return
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("Request %d: %d %s\n", id, resp.StatusCode, string(body))
		}(i)
	}

	wg.Wait()

	fmt.Println("Sending repeated request after completion")

	req, err := http.NewRequest(http.MethodPost, server.URL, bytes.NewBuffer(nil))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Idempotency-Key", key)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	fmt.Printf("Final repeated request: %d %s\n", resp.StatusCode, string(body))
	fmt.Println("Business logic executed:", atomic.LoadInt32(&businessCalls), "time")
}
