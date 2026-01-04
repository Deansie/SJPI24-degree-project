package cmd

import (
	"fmt"

	"github.com/Deansie/SJPI24-degree-project/internal/dns"
	"github.com/spf13/cobra"
)

var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "Sync or manage DNS records in Cloudflare for Kubernetes services/Ingress",
	Long: appLogo + `Idempotently creates, updates, or checks DNS records (e.g., A/CNAME) in Cloudflare to point to your Ingress
load balancer. Handles cases where records already exist (update) or need creation.

Used in deployment flows to ensure domains resolve correctly before applying manifests.

Examples:
  k8s-deploy dns sync --domain app.example.com --target 192.0.2.1
  k8s-deploy dns create --zone example.com --name app --type A --content 192.0.2.1
  k8s-deploy dns list --zone example.com`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := dns.NewClient()
		if err != nil {
			return err
		}
		fmt.Println("✓ Cloudflare authentication successful")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(dnsCmd)
}
