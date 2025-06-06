package bootstrap

import (
	"fmt"
	"goravel/app/models"
	"goravel/config"
	"goravel/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewApp() *fiber.App {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Could not load .env file")
	}

	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())

	// Buat koneksi DB di sini
	db := ConnectDB()

	// Daftarkan Rute dan berikan (suntikkan) koneksi DB
	routes.RegisterApiRoutes(app, db)
	routes.RegisterWebsocketRoutes(app) // Asumsi rute websocket tidak butuh DB, jika butuh, ubah juga

	return app
}

func ConnectDB() *gorm.DB {
	dbConfig := config.LoadDBConfig()

	// UBAH FORMAT DSN UNTUK MYSQL
	// format: "user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
	)

	// GANTI postgres.Open MENJADI mysql.Open
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return db
}

// MigrateDB menjalankan GORM AutoMigrate
func MigrateDB() {
	db := ConnectDB()
	log.Println("Running database migrations...")
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Could not migrate database: %v", err)
	}
	log.Println("Database migration completed successfully.")
}
