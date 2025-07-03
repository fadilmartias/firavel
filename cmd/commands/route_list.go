package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fadilmartias/firavel/bootstrap"

	"github.com/spf13/cobra"
)

var routeListCmd = &cobra.Command{
	Use:   "route:list",
	Short: "List all registered routes",
	Run: func(cmd *cobra.Command, args []string) {
		app := bootstrap.NewApp()

		routes := app.GetRoutes(true) // Gunakan true untuk menyertakan rute internal
		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.TabIndent)

		fmt.Fprintln(writer, "METHOD\tPATH\tNAME") // Ubah Handler menjadi NAME
		fmt.Fprintln(writer, "------\t----\t----")

		for _, route := range routes {
			// Ganti route.Handler dengan route.Name
			// Ini adalah cara yang benar untuk Fiber v2
			fmt.Fprintf(writer, "%s\t%s\t%s\n", route.Method, route.Path, route.Name)
		}

		writer.Flush()
	},
}
