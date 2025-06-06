package commands

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "artisan",
	Short: "Artisan is a CLI tool for Goavel",
	Long:  `A powerful CLI tool inspired by Laravel's Artisan to help manage and build your Goavel application.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Tambahkan semua perintah di sini
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(dbCmd)
	rootCmd.AddCommand(makeCmd)
	rootCmd.AddCommand(routeListCmd)
}