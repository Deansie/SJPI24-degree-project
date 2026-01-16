package dns

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/cloudflare/cloudflare-go"
)

func retryWithBackoff(
	ctx context.Context,
	maxRetries int,
	baseDelay time.Duration,
	fn func() error,
) error {
	var err error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err = fn()
		if err == nil {
			return nil // Success
		}

		if !isRetryable(err) {
			return err // Fail immediately for non-transient errors
		}

		if attempt == maxRetries {
			return err // Exhausted retries
		}

		backoff := float64(baseDelay) * math.Pow(2, float64(attempt))
		jitter := backoff * (0.5 + rand.Float64())
		wait := time.Duration(jitter)

		timer := time.NewTimer(wait)
		select {
		case <-timer.C:
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		}
	}
	return err
}

func isRetryable(err error) bool {
	var cfErr *cloudflare.Error
	if errors.As(err, &cfErr) {
		status := cfErr.StatusCode
		if status == http.StatusTooManyRequests || (status >= 500 && status <= 599) {
			return true
		}
		return false
	}

	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, http.ErrServerClosed) {
		return true
	}

	return false
}
