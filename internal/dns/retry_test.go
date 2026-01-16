package dns

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/cloudflare/cloudflare-go"
)

func TestRetryWithBackoff_SuccessImmediate(t *testing.T) {
	attempts := 0
	err := retryWithBackoff(context.Background(), 5, 10*time.Millisecond, func() error {
		attempts++
		return nil // Success on first try
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if attempts != 1 {
		t.Errorf("expected 1 attempt, got %d", attempts)
	} else {
		t.Logf("PASS: Success on first attempt (attempts: %d)", attempts)
	}
}

func TestRetryWithBackoff_RetriesThenSuccess(t *testing.T) {
	attempts := 0
	maxAttempts := 3
	err := retryWithBackoff(context.Background(), 5, 10*time.Millisecond, func() error {
		attempts++
		if attempts < maxAttempts {
			return &cloudflare.Error{StatusCode: http.StatusTooManyRequests}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected no error after retries, got %v", err)
	}
	if attempts != maxAttempts {
		t.Errorf("expected %d attempts, got %d", maxAttempts, attempts)
	} else {
		t.Logf("PASS: Succeeded after %d attempts (including %d retries)", attempts, attempts-1)
	}
}

func TestRetryWithBackoff_NonRetryableError(t *testing.T) {
	attempts := 0
	err := retryWithBackoff(context.Background(), 5, 10*time.Millisecond, func() error {
		attempts++
		return &cloudflare.Error{StatusCode: http.StatusForbidden}
	})
	if err == nil {
		t.Fatal("expected error for non-retryable case")
	}
	if attempts != 1 {
		t.Errorf("expected only 1 attempt for non-retryable error, got %d", attempts)
	} else {
		t.Logf("PASS: Non-retryable error failed immediately (attempts: %d)", attempts)
	}
}

func TestRetryWithBackoff_ExhaustRetries(t *testing.T) {
	attempts := 0
	err := retryWithBackoff(context.Background(), 2, 10*time.Millisecond, func() error {
		attempts++
		return &cloudflare.Error{StatusCode: 503}
	})
	if err == nil {
		t.Fatal("expected error after exhausting retries")
	}
	expectedAttempts := 3 // Initial + 2 retries
	if attempts != expectedAttempts {
		t.Errorf("expected %d attempts, got %d", expectedAttempts, attempts)
	} else {
		t.Logf("PASS: Exhausted retries after %d attempts", attempts)
	}
}

func TestRetryWithBackoff_ContextTimeout(t *testing.T) {
	attempts := 0
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := retryWithBackoff(ctx, 5, 100*time.Millisecond, func() error {
		attempts++
		time.Sleep(200 * time.Millisecond) // Simulate slow call
		return &cloudflare.Error{StatusCode: 503}
	})
	if err == nil || err != context.DeadlineExceeded {
		t.Fatalf("expected context deadline exceeded, got %v", err)
	}
	if attempts < 1 {
		t.Errorf("expected at least 1 attempt, got %d", attempts)
	} else {
		t.Logf("PASS: Context cancelled after %d attempt(s)", attempts)
	}
}
