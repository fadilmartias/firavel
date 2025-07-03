package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fadilmartias/firavel/app/http/middlewares"
	"github.com/fadilmartias/firavel/app/logger"
	"github.com/fadilmartias/firavel/app/models"
	"github.com/fadilmartias/firavel/app/utils"
	"github.com/fadilmartias/firavel/cmd/cronjob"
	"github.com/fadilmartias/firavel/config"
	"github.com/fadilmartias/firavel/routes"

	"github.com/bytedance/sonic"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	fLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewApp() *fiber.App {
	logger.Init()
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Could not load .env file")
	}

	app := fiber.New(fiber.Config{
		AppName:     os.Getenv("APP_NAME"),
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's a *fiber.Error
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			message := err.Error()
			if message == "" {
				message = "Internal Server Error"
			}

			return utils.ErrorResponse(ctx, utils.ErrorResponseFormat{
				Code:    code,
				Message: message,
				Details: err,
			})
		},
	})
	app.Use(middlewares.LoggerMiddleware())
	app.Use(fLogger.New())
	// Buat koneksi DB di sini
	db := ConnectDB()
	redis := ConnectRedis()
	app.Use(recover.New(recover.Config{
		EnableStackTrace: os.Getenv("APP_ENV") != "production",
	}))
	app.Use(etag.New())
	app.Use(helmet.New())
	app.Use(idempotency.New())
	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        20,
		Expiration: 30 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for")
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
		LimiterMiddleware: limiter.SlidingWindow{},
	}))
	app.Use(cors.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))
	app.Static("/", "./public") // Static file
	cronjob.StartCronJob()

	// Daftarkan Rute dan berikan (suntikkan) koneksi DB
	routes.RegisterApiRoutes(app, db, redis)
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

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return db
}

func ConnectRedis() *redis.Client {
	redisConfig := config.LoadRedisConfig()
	client := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	return client
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
