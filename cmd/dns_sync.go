package cmd

import (
	"context"
	"fmt"

	"github.com/Deansie/SJPI24-degree-project/internal/dns"
	"github.com/Deansie/SJPI24-degree-project/internal/logging"
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

		logging.L().Info("Starting DNS sync", "domain", syncDomain)

		if syncDomain == "" {
			return fmt.Errorf("--domain is required")
		}
		if syncTarget == "" {
			return fmt.Errorf("--target is required")
		}

		client, err := dns.NewClient()
		if err != nil {
			logging.L().Error("Failed to create Cloudflare client", "error", err)
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
			logging.L().Error("DNS sync failed", "error", err)
			return err
		}

		fmt.Println(result.Message)
		logging.L().Info("DNS sync completed", "result", result.Message)
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
