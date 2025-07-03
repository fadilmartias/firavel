package commands

import (
	"fmt"

	"github.com/fadilmartias/firavel/app/models"
	"github.com/fadilmartias/firavel/bootstrap"
	"github.com/fadilmartias/firavel/database/migrations"
	"github.com/fadilmartias/firavel/database/seeders"

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
		db.AutoMigrate(&models.SchemaMigration{})

		seed, _ := cmd.Flags().GetBool("seed")

		for _, migration := range migrations.GetMigrations() {
			if !migrations.HasMigrationRun(db, migration.Name) {
				fmt.Printf("Running migration: %s\n", migration.Name)
				migration.Up(db)
				migrations.RecordMigration(db, migration.Name)
			}
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
		db.AutoMigrate(&models.SchemaMigration{})

		var last models.SchemaMigration
		db.Order("migrated_at DESC").First(&last)

		for _, migration := range migrations.GetMigrations() {
			if migration.Name == last.Name {
				fmt.Printf("Rolling back migration: %s\n", migration.Name)
				migration.Down(db)
				migrations.DeleteMigrationRecord(db, last.Name)
				break
			}
		}
	},
}

var dbMigrateFreshCmd = &cobra.Command{
	Use:   "migrate:fresh",
	Short: "Drop all tables and re-run all migrations",
	Run: func(cmd *cobra.Command, args []string) {
		db := bootstrap.ConnectDB()
		db.AutoMigrate(&models.SchemaMigration{})

		seed, _ := cmd.Flags().GetBool("seed")

		for _, migration := range migrations.GetMigrations() {
			migration.Down(db)
		}
		db.Exec("DELETE FROM schema_migrations")

		for _, migration := range migrations.GetMigrations() {
			migration.Up(db)
			migrations.RecordMigration(db, migration.Name)
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

var dbMigrateStatusCmd = &cobra.Command{
	Use:   "migrate:status",
	Short: "List all migrations and their status",
	Run: func(cmd *cobra.Command, args []string) {
		db := bootstrap.ConnectDB()
		_ = db.AutoMigrate(&models.SchemaMigration{}) // pastikan table schema_migrations tersedia

		applied := map[string]bool{}
		var appliedMigrations []models.SchemaMigration
		db.Find(&appliedMigrations)

		for _, m := range appliedMigrations {
			applied[m.Name] = true
		}

		fmt.Println("\nMIGRATION STATUS")
		fmt.Println("──────────────────────────────")
		fmt.Println(" Status  | Migration Name")
		fmt.Println("─────────┼──────────────────────────────")

		for _, migration := range migrations.GetMigrations() {
			status := "Pending"
			if applied[migration.Name] {
				status = "✅ Ran"
			}
			fmt.Printf(" %-7s | %s\n", status, migration.Name)
		}
		fmt.Println("──────────────────────────────")
	},
}

func init() {
	dbCmd.AddCommand(dbMigrateCmd)
	dbCmd.AddCommand(dbRollbackCmd)
	dbCmd.AddCommand(dbMigrateFreshCmd)
	dbCmd.AddCommand(dbSeedCmd)
	dbCmd.AddCommand(dbMigrateStatusCmd)
	dbMigrateCmd.Flags().BoolP("seed", "s", false, "Seed the database after migration")
}
