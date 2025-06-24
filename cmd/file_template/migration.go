package file_template

const MigrationTemplate = `package migrations

import (
	"log"

	"github.com/fadilmartias/firavel/app/models"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration("{{.MigrationName}}", Up_{{.Timestamp}}_{{.MigrationName}}, Down_{{.Timestamp}}_{{.MigrationName}})
}

func Up_{{.Timestamp}}_{{.MigrationName}}(db *gorm.DB) {
	log.Println("Running migration: {{.MigrationName}} (UP)")
	// TODO: Replace with actual migration logic
	// Example: db.AutoMigrate(&YourModel{})
	err := db.AutoMigrate(&models.{{.Name}}{})
	if err != nil {
		log.Fatalf("Could not create {{.LowerName}} table: %v", err)
	}
	log.Println("Migration {{.MigrationName}} completed successfully.")
}

func Down_{{.Timestamp}}_{{.MigrationName}}(db *gorm.DB) {
	log.Println("Running migration: {{.MigrationName}} (DOWN)")
	// TODO: Replace with actual rollback logic
	// Example: db.Migrator().DropTable("your_table_name")
	err := db.Migrator().DropTable(&models.{{.Name}}{})
	if err != nil {
		log.Fatalf("Could not drop {{.LowerName}} table: %v", err)
	}
	log.Println("Rollback {{.MigrationName}} completed successfully.")
}
`
