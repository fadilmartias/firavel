package main

import (
	"goravel/bootstrap"
	"goravel/config"
	"log"
)

func main() {
	// Buat instance aplikasi dari bootstrap
	app := bootstrap.NewApp()

	// Muat konfigurasi port dari config
	appConfig := config.LoadAppConfig()

	// Jalankan server
	log.Fatal(app.Listen(appConfig.Port))
}
