package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"vchan.in/heyecho/cmd/cli"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "heyecho",
		Short: "HeyEcho CLI",
		Long:  `HeyEcho is a command-line interface for managing the HeyEcho application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Welcome to HeyEcho CLI! Use 'heyecho help' for more information.")
		},
	}

	rootCmd.PersistentFlags().BoolP("help", "h", false, "Display help for command")

	// Add subcommands
	rootCmd.AddCommand(cli.StartCmd())

	// Execute CLI
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
