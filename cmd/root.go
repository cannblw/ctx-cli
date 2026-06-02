package cmd

import "github.com/spf13/cobra"

// NewRootCmd creates the base command with no subcommands wired yet.
// Subcommands are added in main.go via rootCmd.AddCommand(...).
func NewRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ctx",
		Short: "Harness context switching across workstreams",
		Long: `ctx helps you save, list, and switch between workstreams so you can
stay focused when juggling agentic flows, multiple projects, or anything
that pulls your attention in different directions.`,
		// No Run function — root just shows help when invoked without subcommands.
	}
}
