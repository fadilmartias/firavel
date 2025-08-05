package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fadilmartias/firavel/app/jobs"
	"github.com/fadilmartias/firavel/bootstrap"
	"github.com/fadilmartias/firavel/config"
)

func main() {
	// Buat instance aplikasi dari bootstrap
	app, db, redis := bootstrap.NewApp()
	jobs.InitQueue(redis)
	go func() {
		if err := jobs.AsynqServer.Start(jobs.NewHandler(db)); err != nil {
			log.Fatalf("Asynq server error: %v", err)
		}
	}()

	// Muat konfigurasi port dari config
	appConfig := config.LoadAppConfig()

	// Tambahkan handler shutdown di goroutine
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("Shutting down...")
		jobs.AsynqClient.Close()
		jobs.AsynqServer.Stop()
		redis.Close()
		if err := app.Shutdown(); err != nil {
			log.Fatalf("Shutdown error: %v", err)
		}
	}()

	// Jalankan server
	log.Fatal(app.Listen(appConfig.Port))
}
