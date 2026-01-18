package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docsDir string

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Generate Markdown docs for the CLI",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := doc.GenMarkdownTree(rootCmd, docsDir); err != nil {
			return fmt.Errorf("failed to generate docs: %w", err)
		}
		fmt.Printf("Docs generated in %s\n", docsDir)
		return nil
	},
}

func init() {
	docsCmd.Flags().StringVar(&docsDir, "dir", "./docs", "Output directory for docs")
	rootCmd.AddCommand(docsCmd)
}
