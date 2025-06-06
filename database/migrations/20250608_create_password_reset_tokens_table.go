package migrations

import (
	"goravel/app/models"
	"log"

	"gorm.io/gorm"
)

func init() {
	RegisterMigration("create_password_reset_tokens_table", Up_20250608000000_create_password_reset_tokens_table, Down_20250608000000_create_password_reset_tokens_table)
}

// Up_20250608000000_create_password_reset_tokens_table menjalankan migrasi untuk membuat tabel password_reset_tokens
func Up_20250608000000_create_password_reset_tokens_table(db *gorm.DB) {
	log.Println("Running migration: create_password_reset_tokens_table (UP)")
	// Menggunakan AutoMigrate untuk membuat tabel berdasarkan struct model
	// Ini adalah cara GORM yang paling umum untuk "CREATE TABLE"
	err := db.AutoMigrate(&models.PasswordResetToken{})
	if err != nil {
		log.Fatalf("Could not migrate password_reset_tokens table: %v", err)
	}
	log.Println("Migration create_password_reset_tokens_table completed successfully.")
}

// Down_20250608000000_create_password_reset_tokens_table menjalankan rollback migrasi
func Down_20250608000000_create_password_reset_tokens_table(db *gorm.DB) {
	log.Println("Running migration: create_password_reset_tokens_table (DOWN)")
	// Migrator().DropTable akan menghapus tabel
	err := db.Migrator().DropTable(&models.PasswordResetToken{})
	if err != nil {
		log.Fatalf("Could not drop password_reset_tokens table: %v", err)
	}
	log.Println("Rollback create_password_reset_tokens_table completed successfully.")
}
