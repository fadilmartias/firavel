package models

import "time"

type SchemaMigration struct {
	Name       string    `gorm:"primaryKey;size:255"`
	MigratedAt time.Time `gorm:"autoCreateTime;no"`
}

func (SchemaMigration) TableName() string {
	return "schema_migrations"
}
