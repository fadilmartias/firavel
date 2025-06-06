package commands

import (
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
		// Dapatkan koneksi DB
		db := bootstrap.ConnectDB()

		// Panggil migrasi spesifik di sini
		// Dalam sistem nyata, Anda akan memiliki loop atau registri
		// untuk menjalankan semua migrasi yang belum dijalankan.
		migrations.Up_20250607000000_create_users_table(db)
		// migrations.Up_20250608000000_create_products_table(db) // Contoh migrasi lain
	},
}

// Perintah untuk rollback
var dbRollbackCmd = &cobra.Command{
	Use:   "migrate:rollback",
	Short: "Rollback the last database migration",
	Run: func(cmd *cobra.Command, args []string) {
		db := bootstrap.ConnectDB()

		// Panggil fungsi Down dari migrasi yang ingin di-rollback
		migrations.Down_20250607000000_create_users_table(db)
	},
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
	dbCmd.AddCommand(dbRollbackCmd) // Tambahkan perintah rollback
	dbCmd.AddCommand(dbSeedCmd)
}
