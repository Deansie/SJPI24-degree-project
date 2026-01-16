package dns

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Deansie/SJPI24-degree-project/internal/logging"
	"github.com/cloudflare/cloudflare-go"
)

type SyncParams struct {
	Domain string
	Target string
	Type   string // A, CNAME
	Proxy  bool
}

type SyncResult struct {
	Action  string
	Message string
}

func SyncDNSRecord(
	ctx context.Context,
	client *cloudflare.API,
	params SyncParams,
) (*SyncResult, error) {

	// Normalize input
	domain := strings.ToLower(strings.TrimSuffix(params.Domain, "."))
	if domain == "" {
		return nil, fmt.Errorf("domain is required")
	}

	recordType := strings.ToUpper(params.Type)
	if recordType == "" {
		recordType = "A"
	}

	if recordType == "CNAME" && !strings.Contains(params.Target, ".") {
		return nil, fmt.Errorf("CNAME record must point to a hostname")
	}

	// Resolve zone
	zoneID, zoneName, err := findZone(ctx, client, domain)
	if err != nil {
		return nil, err
	}

	if !strings.HasSuffix(domain+".", zoneName+".") {
		return nil, fmt.Errorf("domain %s does not belong to zone %s", domain, zoneName)
	}

	rc := cloudflare.ZoneIdentifier(zoneID)

	var records []cloudflare.DNSRecord
	err = retryWithBackoff(ctx, 5, time.Second, func() error {
		var innerErr error
		records, _, innerErr = client.ListDNSRecords(ctx, rc, cloudflare.ListDNSRecordsParams{
			Name: domain,
			Type: recordType,
		})
		return innerErr
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list DNS records: %w", err)
	}

	// Desired proxy state (only if allowed)
	var desiredProxy *bool
	if proxyAllowed(recordType) {
		p := params.Proxy
		desiredProxy = &p
	}

	switch len(records) {

	case 0:
		// CREATE
		err = retryWithBackoff(ctx, 5, time.Second, func() error {
			var innerErr error
			_, innerErr = client.CreateDNSRecord(ctx, rc, cloudflare.CreateDNSRecordParams{
				Type:    recordType,
				Name:    domain,
				Content: params.Target,
				TTL:     1, // AUTO
				Proxied: desiredProxy,
			})
			return innerErr
		})
		if err != nil {
			logging.L().Error(
				"failed to create DNS record",
				"domain", domain,
				"type", recordType,
				"target", params.Target,
				"proxy", params.Proxy,
				"err", err,
			)
			return nil, fmt.Errorf("failed to create DNS record: %w", err)
		}

		return &SyncResult{
			Action: "created",
			Message: fmt.Sprintf(
				"DNS record created: %s %s -> %s (proxy=%v)",
				recordType,
				domain,
				params.Target,
				params.Proxy,
			),
		}, nil

	case 1:
		record := records[0]

		// Determine current proxy state
		currentProxy := false
		if record.Proxied != nil {
			currentProxy = *record.Proxied
		}

		proxyMatches := true
		if proxyAllowed(recordType) {
			proxyMatches = currentProxy == params.Proxy
		}

		// NOOP only if BOTH content and proxy match
		if record.Content == params.Target && proxyMatches {
			return &SyncResult{
				Action: "noop",
				Message: fmt.Sprintf(
					"DNS record already up to date: %s %s -> %s (proxy=%v)",
					recordType,
					domain,
					params.Target,
					params.Proxy,
				),
			}, nil
		}

		// UPDATE
		err = retryWithBackoff(ctx, 5, time.Second, func() error {
			var innerErr error
			_, innerErr = client.UpdateDNSRecord(ctx, rc, cloudflare.UpdateDNSRecordParams{
				ID:      record.ID,
				Type:    recordType,
				Name:    domain,
				Content: params.Target,
				TTL:     1,            // AUTO
				Proxied: desiredProxy, // nil preserves existing when not allowed
			})
			return innerErr
		})
		if err != nil {
			logging.L().Error(
				"failed to update DNS record",
				"domain", domain,
				"type", recordType,
				"target", params.Target,
				"proxy", params.Proxy,
				"err", err,
			)
			return nil, fmt.Errorf("failed to update DNS record: %w", err)
		}

		return &SyncResult{
			Action: "updated",
			Message: fmt.Sprintf(
				"DNS record updated: %s %s -> %s (proxy=%v)",
				recordType,
				domain,
				params.Target,
				params.Proxy,
			),
		}, nil

	default:
		return nil, fmt.Errorf(
			"multiple DNS records found for %s %s (found %d)",
			domain,
			recordType,
			len(records),
		)
	}
}

func proxyAllowed(recordType string) bool {
	switch recordType {
	case "A", "AAAA", "CNAME":
		return true
	default:
		return false
	}
}

func findZone(
	ctx context.Context,
	client *cloudflare.API,
	domain string,
) (zoneID, zoneName string, err error) {

	labels := strings.Split(domain, ".")
	if len(labels) < 2 {
		return "", "", fmt.Errorf("invalid domain: %s", domain)
	}

	for i := len(labels); i >= 2; i-- {
		candidate := strings.Join(labels[len(labels)-i:], ".")

		var zones []cloudflare.Zone
		err = retryWithBackoff(ctx, 5, time.Second, func() error {
			var innerErr error
			zones, innerErr = client.ListZones(ctx, candidate)
			return innerErr
		})
		if err != nil {
			logging.L().Error(
				"failed to resolve zone for domain",
				"domain", domain,
				"candidate", candidate,
				"err", err,
			)
			return "", "", fmt.Errorf("failed to list zones: %w", err)
		}

		if len(zones) > 0 {
			return zones[0].ID, zones[0].Name, nil
		}
	}

	return "", "", fmt.Errorf("no zone found for domain %s", domain)
}
