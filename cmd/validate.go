package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate Kubernetes YAML manifests against schema and SRE rules",
	Long: `Validates one or more YAML files containing Kubernetes resources (Deployments, Services, Ingress, etc.) 
using official schemas and custom SRE-inspired checks (e.g., resource limits, labels, security contexts).

Fails the deployment pipeline early if issues are found, preventing faulty applies to the cluster.

Examples:
  k8s-deploy validate --file myapp/manifests/*.yaml
  k8s-deploy validate deployment.yaml service.yaml ingress.yaml
  k8s-deploy validate --strict  # Enforce additional best-practice rules`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("validate called")
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
