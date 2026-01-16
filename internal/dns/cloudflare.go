package dns

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Deansie/SJPI24-degree-project/internal/logging"
	"github.com/cloudflare/cloudflare-go"
)

func NewClient() (*cloudflare.API, error) {
	token := os.Getenv("CLOUDFLARE_API_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("CLOUDFLARE_API_TOKEN is not set")
	}

	client, err := cloudflare.NewWithAPIToken(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Cloudflare client: %w", err)
	}

	// Verify token by listing zones (safe, read-only)
	var zones []cloudflare.Zone
	err = retryWithBackoff(context.Background(), 5, time.Second, func() error {
		var innerErr error
		zones, innerErr = client.ListZones(context.Background())
		return innerErr
	})

	if err != nil {
		logging.L().Error("Cloudflare client authentication failed",
			"err", err,
		)
		return nil, fmt.Errorf("Cloudflare authentication failed: %w", err)
	}

	if len(zones) == 0 {
		return nil, fmt.Errorf("Cloudflare authentication failed: no accessible zones")
	}

	return client, nil
}
