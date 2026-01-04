package dns

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/option"
	"github.com/cloudflare/cloudflare-go/v6/zones"
)

func NewClient() (*cloudflare.Client, error) {
	token := os.Getenv("CLOUDFLARE_API_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("CLOUDFLARE_API_TOKEN is not set")
	}

	client := cloudflare.NewClient(
		option.WithAPIToken(token),
	)

	// Verify token by listing zones (zone-scoped, token-safe)
	z, err := client.Zones.List(context.Background(), zones.ZoneListParams{})
	if err != nil {
		return nil, fmt.Errorf("Cloudflare authentication failed: %w", err)
	}

	if len(z.Result) == 0 {
		return nil, fmt.Errorf("Cloudflare authentication failed: no accessible zones")
	}

	return client, nil
}
