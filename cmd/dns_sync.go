package cmd

import (
	"context"
	"fmt"

	"github.com/Deansie/SJPI24-degree-project/internal/dns"
	"github.com/spf13/cobra"
)

var (
	syncDomain string
	syncTarget string
	syncType   string
	syncProxy  bool
)

var dnsSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Idempotently create or update DNS records in Cloudflare",
	Long: `Synchronize a DNS record in Cloudflare.

The sync command ensures that a DNS record exists and matches the desired
state. It can be safely executed multiple times without creating duplicates,
making it suitable for CI/CD pipelines.

If the record does not exist, it will be created.
If the record exists but differs, it will be updated.
If the record already matches, no changes are made.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if syncDomain == "" {
			return fmt.Errorf("--domain is required")
		}
		if syncTarget == "" {
			return fmt.Errorf("--target is required")
		}

		client, err := dns.NewClient()
		if err != nil {
			return err
		}

		result, err := dns.SyncDNSRecord(
			context.Background(),
			client,
			dns.SyncParams{
				Domain: syncDomain,
				Target: syncTarget,
				Type:   syncType,
				Proxy:  syncProxy,
			},
		)
		if err != nil {
			return err
		}

		fmt.Println(result.Message)
		return nil
	},
}

func init() {
	dnsSyncCmd.Flags().StringVar(&syncDomain, "domain", "", "Fully qualified domain name (e.g. app.example.com)")
	dnsSyncCmd.Flags().StringVar(&syncTarget, "target", "", "IP address or hostname the record should point to")
	dnsSyncCmd.Flags().StringVar(&syncType, "type", "A", "DNS record type (A or CNAME)")
	dnsSyncCmd.Flags().BoolVar(&syncProxy, "proxy", false, "Enable Cloudflare proxy (orange cloud)")

	dnsCmd.AddCommand(dnsSyncCmd)
}
