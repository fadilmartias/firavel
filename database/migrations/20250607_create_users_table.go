package migrations

import (
	"goravel/app/models"
	"log"

	"gorm.io/gorm"
)

func init() {
	RegisterMigration("create_users_table", Up_20250607000000_create_users_table, Down_20250607000000_create_users_table)
}

// Up_20250607000000_create_users_table menjalankan migrasi untuk membuat tabel users
func Up_20250607000000_create_users_table(db *gorm.DB) {
	log.Println("Running migration: create_users_table (UP)")
	// Menggunakan AutoMigrate untuk membuat tabel berdasarkan struct model
	// Ini adalah cara GORM yang paling umum untuk "CREATE TABLE"
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Could not migrate user table: %v", err)
	}
	log.Println("Migration create_users_table completed successfully.")
}

// Down_20250607000000_create_users_table menjalankan rollback migrasi
func Down_20250607000000_create_users_table(db *gorm.DB) {
	log.Println("Running migration: create_users_table (DOWN)")
	// Migrator().DropTable akan menghapus tabel
	err := db.Migrator().DropTable(&models.User{})
	if err != nil {
		log.Fatalf("Could not drop user table: %v", err)
	}
	log.Println("Rollback create_users_table completed successfully.")
}
