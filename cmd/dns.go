package cmd

import (
	"fmt"

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
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Dns called")
	},
}

func init() {
	rootCmd.AddCommand(dnsCmd)
}
