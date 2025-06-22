package seeders

import (
	"log"

	"github.com/fadilmartias/firavel/bootstrap"
)

// RunAllSeeders adalah entry point untuk menjalankan semua seeder.
func RunAllSeeders() {
	db := bootstrap.ConnectDB()
	log.Println("Running all seeders...")

	SeedUsers(db, 20) // Buat 20 user palsu

	log.Println("All seeders completed.")
}
