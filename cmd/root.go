package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/Deansie/SJPI24-degree-project/internal/logging"
	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var version = "dev"

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

var verbose bool

var rootCmd = &cobra.Command{
	Use:   "k8s-deploy",
	Short: "A CLI tool for validating Kubernetes manifests, automating Ingress with TLS, and syncing DNS records.",
	Long: appLogo + `k8s-deploy simplifies and secures Kubernetes deployments by validating YAML manifests against strict
rules, automating Ingress resources with cert-manager for Let's Encrypt TLS, and managing DNS records 
in Cloudflare. Designed for SRE best practices, it helps prevent deployment errors, ensures end-to-end
encryption,and streamlines CI/CD pipelines.`,
	Version: "1.0",
}

func Execute() {
	if isInteractive() && shouldClearScreen() {
		clearScreen()
	}

	if err := fang.Execute(
		context.Background(),
		rootCmd,
		fang.WithVersion(version), // Set version during build 'go build -ldflags "-X github.com/Deansie/SJPI24-degree-project/cmd.version=0.3.0" -o k8s-deploy'
	); err != nil {
		os.Exit(1)
	}
}

func clearScreen() {
	fmt.Print("\033c")
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose (debug) logging")
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		logging.SetVerbose(verbose)
	}
}

var isInteractive = func() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

func shouldClearScreen() bool {
	if len(os.Args) == 1 {
		return true
	}

	for _, arg := range os.Args[1:] {
		if arg == "--help" || arg == "-h" || arg == "help" {
			return true
		}
	}

	return false
}
