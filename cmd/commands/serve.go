package commands

import (
	"log"

	"github.com/fadilmartias/firavel/bootstrap"
	"github.com/fadilmartias/firavel/config"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the web server",
	Run: func(cmd *cobra.Command, args []string) {
		app := bootstrap.NewApp()
		appConfig := config.LoadAppConfig()

		log.Printf("Starting server on http://localhost%s", appConfig.Port)
		log.Fatal(app.Listen(appConfig.Port))
	},
}
