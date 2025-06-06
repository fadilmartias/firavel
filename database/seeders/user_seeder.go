package seeders

import (
	"fmt"
	_ "goravel/app/models"
	"goravel/database/factories"
	"log"

	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB, count int) {
	log.Printf("Seeding %d users...", count)
	for i := 0; i < count; i++ {
		user := factories.NewUser()
		if err := user.HashPassword(user.Password); err != nil {
			log.Printf("Failed to hash password for user %s: %v", user.Email, err)
			continue
		}

		result := db.Create(&user)
		if result.Error != nil {
			log.Printf("Could not seed user: %v", result.Error)
		}
	}
	fmt.Printf("Seeded %d users successfully.\n", count)
}
