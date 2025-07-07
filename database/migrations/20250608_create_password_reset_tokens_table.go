package migrations

import (
	"log"

	"github.com/fadilmartias/firavel/app/models"

	"gorm.io/gorm"
)

func init() {
	RegisterMigration("create_password_reset_tokens_table", Up_20250608000000_create_password_reset_tokens_table, Down_20250608000000_create_password_reset_tokens_table)
}

// Up_20250608000000_create_password_reset_tokens_table menjalankan migrasi untuk membuat tabel password_reset_tokens
func Up_20250608000000_create_password_reset_tokens_table(db *gorm.DB) {
	log.Println("Running migration: create_password_reset_tokens_table (UP)")

	if !db.Migrator().HasTable(&models.PasswordResetToken{}) {
		err := db.Migrator().CreateTable(&models.PasswordResetToken{})
		if err != nil {
			log.Fatalf("Could not create password_reset_tokens table: %v", err)
		}
	}

	log.Println("Migration create_password_reset_tokens_table completed successfully.")
}

// Down_20250608000000_create_password_reset_tokens_table menjalankan rollback migrasi
func Down_20250608000000_create_password_reset_tokens_table(db *gorm.DB) {
	log.Println("Running migration: create_password_reset_tokens_table (DOWN)")

	if db.Migrator().HasTable(&models.PasswordResetToken{}) {
		err := db.Migrator().DropTable(&models.PasswordResetToken{})
		if err != nil {
			log.Fatalf("Could not drop password_reset_tokens table: %v", err)
		}
	}
	log.Println("Rollback create_password_reset_tokens_table completed successfully.")
}
