package migrations

import "gorm.io/gorm"

type Migration struct {
	Name string
	Up   func(db *gorm.DB)
	Down func(db *gorm.DB)
}

var migrationList []Migration

func RegisterMigration(name string, upFunc func(db *gorm.DB), downFunc func(db *gorm.DB)) {
	migrationList = append(migrationList, Migration{
		Name: name,
		Up:   upFunc,
		Down: downFunc,
	})
}

func GetMigrations() []Migration {
	return migrationList
}
