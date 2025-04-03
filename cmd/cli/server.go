package cli

import (
	"github.com/spf13/cobra"
	"vchan.in/heyecho/cmd/app"
	"vchan.in/heyecho/cmd/seed"
	"vchan.in/heyecho/db"
)

// StartCmd represents the start command
func StartCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "start",
		Short: "Start the server",
		Long:  `This command starts the HeyEcho server.`,
		Run: func(cmd *cobra.Command, args []string) {
			db.InitDB()
			seed.Seed()
			app.Server()
		},
	}

	cmd.PersistentFlags().BoolP("help", "h", false, "Display help for command")

	return cmd
}
