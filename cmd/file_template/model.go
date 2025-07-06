package file_template

const ModelTemplate = `package models

import (
	"time"
)

type {{.Name}} struct {
	ID              string     ` + "`gorm:\"primarykey;size:7\"`" + `
	CreatedAt       time.Time  ` + "`gorm:\"not null\"`" + `
	UpdatedAt       time.Time
}`
