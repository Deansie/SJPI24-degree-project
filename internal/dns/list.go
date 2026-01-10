package dns

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/cloudflare/cloudflare-go"
)

type DNSRecord struct {
	Type    string
	Name    string
	Content string
	TTL     int
}

func ListDNSRecords(
	ctx context.Context,
	client *cloudflare.API,
	zoneName string,
	filterType string,
) ([]DNSRecord, error) {

	zones, err := client.ListZones(ctx, zoneName)
	if err != nil {
		return nil, fmt.Errorf("failed to find zone: %w", err)
	}
	if len(zones) == 0 {
		return nil, fmt.Errorf("zone not found: %s", zoneName)
	}

	rc := cloudflare.ZoneIdentifier(zones[0].ID)

	records, _, err := client.ListDNSRecords(ctx, rc, cloudflare.ListDNSRecordsParams{})
	if err != nil {
		return nil, fmt.Errorf("failed to list DNS records: %w", err)
	}

	filterType = strings.ToUpper(filterType)

	out := make([]DNSRecord, 0, len(records))
	for _, r := range records {
		if filterType != "" && r.Type != filterType {
			continue
		}

		out = append(out, DNSRecord{
			Type:    r.Type,
			Name:    r.Name,
			Content: r.Content,
			TTL:     r.TTL,
		})
	}

	// Sort by TYPE, then NAME
	sort.Slice(out, func(i, j int) bool {
		if out[i].Type == out[j].Type {
			return out[i].Name < out[j].Name
		}
		return out[i].Type < out[j].Type
	})

	return out, nil
}
