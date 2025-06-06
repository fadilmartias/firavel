package seeders

import (
	"goravel/bootstrap"
	"log"
)

// RunAllSeeders adalah entry point untuk menjalankan semua seeder.
func RunAllSeeders() {
	db := bootstrap.ConnectDB()
	log.Println("Running all seeders...")

	// Panggil seeder spesifik di sini
	SeedUsers(db, 20) // Buat 20 user palsu

	log.Println("All seeders completed.")
}
