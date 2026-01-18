package cmd

import (
	"github.com/Deansie/SJPI24-degree-project/internal/logging"
	"github.com/spf13/cobra"
)

var ingressCmd = &cobra.Command{
	Use:   "ingress",
	Short: "Setup and manage Kubernetes Ingress resources with automatic TLS via cert-manager",
	Long: appLogo + `Creates or updates Ingress manifests with cert-manager annotations for automatic Let's Encrypt
certificates. Supports initial setup flows and ensures end-to-end TLS encryption.

Integrates seamlessly with DNS sync for full external access automation.

Examples:
	k8s-deploy ingress setup --name myapp --host app.example.com --service myapp-service --affinity --limit-rps=10
	k8s-deploy ingress setup --name myapp --host app.example.com --service myapp-service --enable-cors --cors-allow-origin=https://frontend.com
/* 	k8s-deploy ingress create --file ingress.yaml  # With auto-annotations added
	k8s-deploy ingress list */
	`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logging.L().Info("Entering Ingress command group")
	},
}

func init() {
	rootCmd.AddCommand(ingressCmd)
	ingressCmd.AddCommand(ingressSetupCmd)
}
