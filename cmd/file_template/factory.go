package file_template

const FactoryTemplate = `package factories

import (
	"log"
	"time"

	"github.com/fadilmartias/firavel/app/models"

	"github.com/bxcodec/faker/v3"
)

func New{{.Name}}() models.{{.Name}} {
	var item models.{{.Name}}
	err := faker.FakeData(&item)
	if err != nil {
		log.Printf("Error faking {{.LowerName}} data: %v", err)
	}
	// Optional: set default values
	// item.Password = "password"
	return item
}
`
