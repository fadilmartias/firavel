package main

import (
	"log"

	"github.com/fadilmartias/firavel/cmd/commands"

	"github.com/joho/godotenv"
)

func main() {
	// CLI juga perlu .env untuk koneksi DB, dll.
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Could not load .env file")
	}

	if err := commands.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
