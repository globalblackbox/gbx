// cmd/root.go
package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gbx",
	Short: "GBX is a tool to interact with Global Blackbox services",
	Long: `GBX allows you to sign up, manage your account,
and interact with Global Blackbox services through a command-line interface.`,
	// You can add a default action here if needed
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	// Add subcommands
	rootCmd.AddCommand(signupCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		// Styled error message
		style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("1"))
		fmt.Fprintf(os.Stderr, "%s: %v\n", style.Render("Error"), err)
		os.Exit(1)
	}
}
