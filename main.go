package main

import (
	"log"

	"github.com/fadilmartias/firavel/bootstrap"
	"github.com/fadilmartias/firavel/config"
)

func main() {
	// Buat instance aplikasi dari bootstrap
	app := bootstrap.NewApp()

	// Muat konfigurasi port dari config
	appConfig := config.LoadAppConfig()

	// Jalankan server
	log.Fatal(app.Listen(appConfig.Port))
}
