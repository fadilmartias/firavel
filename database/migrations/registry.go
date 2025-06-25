// database/migrations/registry.go (atau file lain yang inisialisasi RegisterMigration)
package migrations

import (
	"github.com/fadilmartias/firavel/app/models"
	"gorm.io/gorm"
)

type Migration struct {
	Name string
	Up   func(db *gorm.DB)
	Down func(db *gorm.DB)
}

var migrationList []Migration

func RegisterMigration(name string, up func(*gorm.DB), down func(*gorm.DB)) {
	migrationList = append(migrationList, Migration{
		Name: name,
		Up:   up,
		Down: down,
	})
}

func GetMigrations() []Migration {
	return migrationList
}

func HasMigrationRun(db *gorm.DB, name string) bool {
	var count int64
	db.Model(&models.SchemaMigration{}).Where("name = ?", name).Count(&count)
	return count > 0
}

func RecordMigration(db *gorm.DB, name string) {
	db.Create(&models.SchemaMigration{Name: name})
}

func DeleteMigrationRecord(db *gorm.DB, name string) {
	db.Delete(&models.SchemaMigration{}, "name = ?", name)
}
