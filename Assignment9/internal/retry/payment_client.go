package retry

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net"
	"net/http"
	"time"
)

const (
	BaseDelay = 500 * time.Millisecond
	MaxDelay  = 5 * time.Second
)

type PaymentClient struct {
	Client     *http.Client
	MaxRetries int
}

func IsRetryable(resp *http.Response, err error) bool {
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return false
		}

		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return true
		}

		return true
	}

	if resp == nil {
		return false
	}

	switch resp.StatusCode {
	case 429, 500, 502, 503, 504:
		return true
	case 401, 404:
		return false
	default:
		return false
	}
}

func CalculateBackoff(attempt int) time.Duration {
	backoff := BaseDelay * time.Duration(math.Pow(2, float64(attempt-1)))
	if backoff > MaxDelay {
		backoff = MaxDelay
	}

	if backoff <= 0 {
		return BaseDelay
	}

	return time.Duration(rand.Int63n(int64(backoff))) + 1
}

func (p *PaymentClient) ExecutePayment(ctx context.Context, url string) ([]byte, error) {
	if p.Client == nil {
		p.Client = http.DefaultClient
	}

	for attempt := 1; attempt <= p.MaxRetries; attempt++ {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
		if err != nil {
			return nil, err
		}

		resp, err := p.Client.Do(req)

		if resp != nil && resp.Body != nil {
			body, readErr := io.ReadAll(resp.Body)
			resp.Body.Close()

			if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
				fmt.Printf("Attempt %d: Success!\n", attempt)
				return body, nil
			}

			if readErr != nil {
				return nil, readErr
			}
		}

		if !IsRetryable(resp, err) {
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("non retryable status: %d", resp.StatusCode)
		}

		if attempt == p.MaxRetries {
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("failed after %d attempts, last status: %d", p.MaxRetries, resp.StatusCode)
		}

		wait := CalculateBackoff(attempt)
		status := 0
		if resp != nil {
			status = resp.StatusCode
		}

		fmt.Printf("Attempt %d failed: status %d, waiting %v...\n", attempt, status, wait)

		select {
		case <-time.After(wait):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return nil, fmt.Errorf("payment failed")
}
