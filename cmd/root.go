package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "gbx",
	Short: "gbx is the CLI tool to interact with the Global Blackbox API",
	Long:  `GBX allows you to sign up and interact with Global Blackbox API through a command-line interface.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	rootCmd.AddCommand(signupCmd)
	rootCmd.AddCommand(logsCmd)

	if err := rootCmd.Execute(); err != nil {
		style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("1"))
		fmt.Fprintf(os.Stderr, "%s: %v\n", style.Render("Error"), err)
		os.Exit(1)
	}
}
