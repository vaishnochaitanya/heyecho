package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// helpCmd represents the help command
var HelpCmd = &cobra.Command{
	Use:   "help",
	Short: "Help command",
	Long:  `This command provides help information for the CLI commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Available commands:")
		fmt.Println("  help       Show help information")
		fmt.Println("  start      Start the server")
		fmt.Println("  stop       Stop the server")
		fmt.Println("  status     Restart the server")
		fmt.Println("  backup     Backup the application data")
		fmt.Println("  restore    Restore the application data")
		fmt.Println("  admin      Manage admin users")
	},
}
