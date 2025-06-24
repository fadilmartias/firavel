package commands

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "artisan",
	Short: "Artisan is a CLI tool for Goravel",
	Long:  `A powerful CLI tool inspired by Laravel's Artisan to help manage and build your Goravel application.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Tambahkan semua perintah di sini
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(dbCmd)
	rootCmd.AddCommand(makeFileCmd)
	rootCmd.AddCommand(routeListCmd)
	rootCmd.AddCommand(makeResponseCmd)
}
