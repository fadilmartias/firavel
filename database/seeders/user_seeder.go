package seeders

import (
	"fmt"
	"goravel/app/models"
	_ "goravel/app/models"
	"goravel/database/factories"
	"log"

	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB, count int) {
	log.Printf("Seeding %d users...", count)
	users := make([]models.User, 0, count)
	admin := models.User{
		Id:       "00000A1",
		Name:     "Admin",
		Email:    "admin@gmail.com",
		Password: "namakau123",
	}
	if err := admin.HashPassword(admin.Password); err != nil {
		log.Printf("Failed to hash password for admin: %v", err)
	}
	users = append(users, admin)
	user := models.User{
		Id:       "00000A2",
		Name:     "User",
		Email:    "user@gmail.com",
		Password: "namakau123",
	}
	if err := user.HashPassword(user.Password); err != nil {
		log.Printf("Failed to hash password for user: %v", err)
	}
	users = append(users, user)
	for i := 0; i < count; i++ {
		user := factories.NewUser()
		if err := user.HashPassword(user.Password); err != nil {
			log.Printf("Failed to hash password for user %s: %v", user.Email, err)
			continue
		}
		users = append(users, user)
	}

	result := db.CreateInBatches(&users, 100)
	if result.Error != nil {
		log.Printf("Could not seed user: %v", result.Error)
	}
	fmt.Printf("Seeded %d users successfully.\n", count)
}
