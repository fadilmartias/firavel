package file_template

const ModelTemplate = `package models

import (
	"time"
)

type {{.Name}} struct {
	ID              string     ` + "`gorm:\"primarykey;not null;size:7\" json:\"id,omitempty\"`" + `
	CreatedAt       time.Time  ` + "`gorm:\"type:timestamp;default:CURRENT_TIMESTAMP\" json:\"created_at,omitempty\"`" + `
	UpdatedAt       *time.Time ` + "`gorm:\"type:timestamp;default:NULL ON UPDATE CURRENT_TIMESTAMP\" json:\"updated_at,omitempty\"`" + `
}`
