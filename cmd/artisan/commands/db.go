package commands

import (
	"fmt"
	"goravel/bootstrap"
	"goravel/database/migrations"
	"goravel/database/seeders"

	"github.com/spf13/cobra"
)

// Grup perintah DB
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database related commands",
}

var dbMigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		db := bootstrap.ConnectDB()

		seed, _ := cmd.Flags().GetBool("seed")

		for _, migration := range migrations.GetMigrations() {
			fmt.Printf("Running migration: %s\n", migration.Name)
			migration.Up(db)
		}

		if seed {
			seeders.RunAllSeeders()
		}
	},
}

var dbRollbackCmd = &cobra.Command{
	Use:   "migrate:rollback",
	Short: "Rollback the last database migration",
	Run: func(cmd *cobra.Command, args []string) {
		db := bootstrap.ConnectDB()

		for _, migration := range migrations.GetMigrations() {
			fmt.Printf("Rolling back migration: %s\n", migration.Name)
			migration.Down(db)
		}
	},
}

var dbMigrateFreshCmd = &cobra.Command{
	Use:   "migrate:fresh",
	Short: "Drop all tables and re-run all migrations",
	Run: func(cmd *cobra.Command, args []string) {
		db := bootstrap.ConnectDB()

		seed, _ := cmd.Flags().GetBool("seed")

		for _, migration := range migrations.GetMigrations() {
			fmt.Printf("Dropping table: %s\n", migration.Name)
			migration.Down(db)
		}

		for _, migration := range migrations.GetMigrations() {
			fmt.Printf("Running migration: %s\n", migration.Name)
			migration.Up(db)
		}

		if seed {
			seeders.RunAllSeeders()
		}
	},
}

func init() {
	dbMigrateFreshCmd.Flags().BoolP("seed", "s", false, "Seed the database after migration")
}

var dbSeedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with records",
	Run: func(cmd *cobra.Command, args []string) {
		seeders.RunAllSeeders()
	},
}

func init() {
	dbCmd.AddCommand(dbMigrateCmd)
	dbCmd.AddCommand(dbRollbackCmd)
	dbCmd.AddCommand(dbMigrateFreshCmd) // Tambahkan perintah migrate:fresh
	dbCmd.AddCommand(dbSeedCmd)

	dbMigrateCmd.Flags().BoolP("seed", "s", false, "Seed the database after migration")
}
