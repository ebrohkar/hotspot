// hotspot is a CLI tool for tracking and visualizing code hotspots
// in a Git repository based on change frequency and complexity metrics.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version is set at build time via ldflags
	Version = "dev"
	// Commit is set at build time via ldflags
	Commit = "none"
	// Date is set at build time via ldflags
	Date = "unknown"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "hotspot",
	Short: "Identify hotspots in your Git repository",
	Long: `hotspot analyzes your Git repository to identify files that change
frequently and have high complexity — the most risky areas of your codebase.

By combining change frequency with code complexity metrics, hotspot helps
you prioritize refactoring efforts and technical debt reduction.

See https://github.com/huangsam/hotspot for the upstream project.
This is a personal fork for learning and experimentation.`,
	SilenceUsage:  true,
	SilenceErrors: true, // handle errors manually in main() for cleaner output
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		// Print a compact one-liner; include fork notice so it's clear this
		// binary is not the upstream release.
		fmt.Printf("hotspot %s (commit: %s, built: %s) [personal fork]\n", Version, Commit, Date)
		// Also print upstream URL for easy reference when sharing output with others
		fmt.Println("upstream: https://github.com/huangsam/hotspot")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	// Run version info on bare invocation so I don't have to remember the subcommand
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		versionCmd.Run(cmd, args)
		return nil
	}
}
