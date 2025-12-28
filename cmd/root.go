package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

const appLogo = `

        ░██        ░██████                ░███████                         ░██
        ░██       ░██   ░██               ░██   ░██                        ░██
        ░██    ░██░██   ░██  ░███████     ░██    ░██  ░███████  ░████████  ░██  ░███████  ░██    ░██
        ░██   ░██  ░██████   ░██          ░██    ░██ ░██    ░██ ░██    ░██ ░██ ░██    ░██ ░██    ░██
        ░███████  ░██   ░██  ░███████     ░██    ░██ ░█████████ ░██    ░██ ░██ ░██    ░██ ░██    ░██
        ░██   ░██ ░██   ░██       ░██     ░██   ░██  ░██        ░███   ░██ ░██ ░██    ░██ ░██    ███
        ░██    ░██ ░██████   ░███████     ░███████    ░███████  ░██░█████  ░██  ░███████   ░█████░██
                                                                ░██                              ░██
                                                                ░██                        ░███████
																
------------------------------------------------------------------------------------------------------------

`

var rootCmd = &cobra.Command{
	Use:   "k8s-deploy",
	Short: "A CLI tool for validating Kubernetes manifests, automating Ingress with TLS, and syncing DNS records.",
	Long: appLogo + `k8s-deploy simplifies and secures Kubernetes deployments by validating YAML manifests against strict
rules, automating Ingress resources with cert-manager for Let's Encrypt TLS, and managing DNS records 
in Cloudflare. Designed for SRE best practices, it helps prevent deployment errors, ensures end-to-end
encryption,and streamlines CI/CD pipelines.`,
	Version: "0.3",
}

func Execute() {
	clearScreen()

	if err := fang.Execute(context.Background(), rootCmd); err != nil {
		os.Exit(1)
	}
}

func clearScreen() {
	fmt.Print("\033c")
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
}
