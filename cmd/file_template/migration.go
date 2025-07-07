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

	if !db.Migrator().HasTable(&models.{{.Name}}{}) {
		err := db.Migrator().CreateTable(&models.{{.Name}}{})
		if err != nil {
			log.Fatalf("Could not create {{.LowerName}} table: %v", err)
		}
	}

	log.Println("Migration {{.MigrationName}} completed successfully.")
}

func Down_{{.Timestamp}}_{{.MigrationName}}(db *gorm.DB) {
	log.Println("Running migration: {{.MigrationName}} (DOWN)")

	if db.Migrator().HasTable(&models.{{.Name}}{}) {
		err := db.Migrator().DropTable(&models.{{.Name}}{})
		if err != nil {
			log.Fatalf("Could not drop {{.LowerName}} table: %v", err)
		}
	}
	log.Println("Rollback {{.MigrationName}} completed successfully.")
}
`
