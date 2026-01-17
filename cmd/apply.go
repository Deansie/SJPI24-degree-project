package cmd

import "github.com/spf13/cobra"

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply Kubernetes manifests to the cluster",
	Long: appLogo + `Applies validated YAML manifests to the Kubernetes cluster using client-go.
Supports kubeconfig for authentication and namespace/context overrides.

Examples:
  k8s-deploy apply exec --file manifest.yaml --namespace myapp
  k8s-deploy apply exec --dir manifests/ --dry-run`,
}

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.AddCommand(applyExecCmd) // For subcommand modularity
}
