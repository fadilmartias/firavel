package file_template

const SeederTemplate = `package seeders

import (
	"fmt"
	"log"

	"github.com/fadilmartias/firavel/app/models"

	"gorm.io/gorm"
)

func Seed{{.Name}}(db *gorm.DB, count int) {
	log.Printf("Seeding %d {{.LowerName}}...", count)
	var items []models.{{.Name}}

	// Optional: predefined example
	sample := models.{{.Name}}{
		// TODO: fill default values here
	}
	items = append(items, sample)

	// for i := 0; i < count; i++ {
	// 	item := factories.New{{.Name}}()
	// 	// Optional: modify item before hashing/storing
	// 	items = append(items, item)
	// }

	result := db.CreateInBatches(&items, 100)
	if result.Error != nil {
		log.Printf("Could not seed {{.LowerName}}: %v", result.Error)
	} else {
		fmt.Printf("Seeded %d {{.LowerName}} successfully.\n", count)
	}
}
`
