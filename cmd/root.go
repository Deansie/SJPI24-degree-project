package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "k8s-deploy",
	Short: "A CLI tool for validating Kubernetes manifests, automating Ingress with TLS, and syncing DNS records.",
	Long: `k8s-deploy simplifies and secures Kubernetes deployments by validating YAML manifests against strict rules,
automating Ingress resources with cert-manager for Let's Encrypt TLS, and managing DNS records in Cloudflare.

Designed for SRE best practices, it helps prevent deployment errors, ensures end-to-end encryption, and streamlines CI/CD pipelines.`,
	Version: "0.2",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
}
