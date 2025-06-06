package commands

import (
	"github.com/spf13/cobra"
	"goravel/app/models"
	"goravel/bootstrap"
	"goravel/database/migrations"
	"goravel/database/seeders"
	"log"
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
		// Dapatkan koneksi DB
		db := bootstrap.ConnectDB()

		seed, _ := cmd.Flags().GetBool("seed")
		if seed {
			seeders.RunAllSeeders()
		}

		// Panggil migrasi spesifik di sini
		// Dalam sistem nyata, Anda akan memiliki loop atau registri
		// untuk menjalankan semua migrasi yang belum dijalankan.
		migrations.Up_20250607000000_create_users_table(db)
		// migrations.Up_20250608000000_create_products_table(db) // Contoh migrasi lain
	},
}

var dbRollbackCmd = &cobra.Command{
	Use:   "migrate:rollback",
	Short: "Rollback the last database migration",
	Run: func(cmd *cobra.Command, args []string) {
		db := bootstrap.ConnectDB()

		// Panggil fungsi Down dari migrasi yang ingin di-rollback
		migrations.Down_20250607000000_create_users_table(db)
	},
}

var dbMigrateFreshCmd = &cobra.Command{
	Use:   "migrate:fresh",
	Short: "Drop all tables and re-run all migrations",
	Run: func(cmd *cobra.Command, args []string) {
		db := bootstrap.ConnectDB()

		seed, _ := cmd.Flags().GetBool("seed")

		// Jatuhkan semua tabel
		err := db.Migrator().DropTable(&models.User{})
		if err != nil {
			log.Fatalf("Could not drop table: %v", err)
		}

		// Jalankan semua migrasi lagi
		migrations.Up_20250607000000_create_users_table(db)

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
