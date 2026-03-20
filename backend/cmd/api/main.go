package main

import (
	"context"
	"log"
	"time"

	"PulsePoll/internal/config"
	postgresinfra "PulsePoll/internal/infrastructure/postgres"
	redisinfra "PulsePoll/internal/infrastructure/redis"
	httproutes "PulsePoll/internal/transport/http/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	db, err := postgresinfra.New(cfg.Postgres)
	if err != nil {
		log.Fatalf("postgres init failed: %v", err)
	}

	redisClient, err := redisinfra.New(ctx, cfg.Redis)
	if err != nil {
		log.Fatalf("redis init failed: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("postgres sql db access failed: %v", err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("postgres close error: %v", err)
		}
	}()
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Printf("redis close error: %v", err)
		}
	}()

	app := fiber.New()

	// CORS (Cross-Origin Resource Sharing / Çapraz Kaynak Paylaşımı)
	// Frontend: http://localhost:3000  -> Backend: http://localhost:8080
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://127.0.0.1:3000",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	httproutes.RegisterAuthRoutes(app, db, cfg.Auth)

	app.Get("/health", func(c *fiber.Ctx) error {
		healthCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		postgresOK := sqlDB.PingContext(healthCtx) == nil
		redisOK := redisClient.Ping(healthCtx).Err() == nil

		if postgresOK && redisOK {
			return c.JSON(fiber.Map{
				"ok":    true,
				"redis": "up",
				"db":    "up",
			})
		}

		status := fiber.Map{
			"ok":    false,
			"redis": "down",
			"db":    "down",
		}
		if redisOK {
			status["redis"] = "up"
		}
		if postgresOK {
			status["db"] = "up"
		}

		return c.Status(fiber.StatusServiceUnavailable).JSON(status)
	})

	log.Printf("API listening on :%s", cfg.App.Port)
	log.Fatal(app.Listen(":" + cfg.App.Port))
}
