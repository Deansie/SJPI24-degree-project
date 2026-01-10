package cmd

import (
	"context"
	"fmt"

	"github.com/Deansie/SJPI24-degree-project/internal/dns"
	"github.com/spf13/cobra"
)

var (
	listZone string
	listType string
)

var dnsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List DNS records in a Cloudflare zone",
	Long: `List DNS records for a Cloudflare zone.

This command is read-only and safe for CI/CD pipelines.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if listZone == "" {
			return fmt.Errorf("--zone is required")
		}

		client, err := dns.NewClient()
		if err != nil {
			return err
		}

		records, err := dns.ListDNSRecords(
			context.Background(),
			client,
			listZone,
			listType,
		)
		if err != nil {
			return err
		}

		if len(records) == 0 {
			fmt.Println("No DNS records found.")
			return nil
		}

		fmt.Printf("%-6s %-40s %-55s %-6s\n", "TYPE", "NAME", "CONTENT", "TTL")
		for _, r := range records {
			fmt.Printf(
				"%-6s %-40s %-55s %-6s\n",
				r.Type,
				r.Name,
				trim(r.Content, 55),
				formatTTL(r.TTL),
			)
		}

		return nil
	},
}

func init() {
	dnsListCmd.Flags().StringVar(&listZone, "zone", "", "DNS zone (e.g. example.com)")
	dnsListCmd.Flags().StringVar(&listType, "type", "", "Filter by record type (A, CNAME, MX, TXT)")
	dnsCmd.AddCommand(dnsListCmd)
}

func trim(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func formatTTL(ttl int) string {
	if ttl == 1 {
		return "AUTO"
	}
	return fmt.Sprintf("%d", ttl)
}
