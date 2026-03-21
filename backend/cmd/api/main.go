package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"PulsePoll/internal/config"
	postgresinfra "PulsePoll/internal/infrastructure/postgres"
	redisinfra "PulsePoll/internal/infrastructure/redis"
	httproutes "PulsePoll/internal/transport/http/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	log.Printf("starting API service")
	cfg := config.Load()
	ctx := context.Background()

	log.Printf("connecting postgres host=%s port=%s db=%s", cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DBName)
	db, err := postgresinfra.New(cfg.Postgres)
	if err != nil {
		log.Fatalf("postgres init failed: %v", err)
	}
	log.Printf("postgres connection ready")

	log.Printf("connecting redis addr=%s db=%d", cfg.Redis.Addr, cfg.Redis.DB)
	redisClient, err := redisinfra.New(ctx, cfg.Redis)
	if err != nil {
		log.Fatalf("redis init failed: %v", err)
	}
	log.Printf("redis connection ready")

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("postgres sql db access failed: %v", err)
	}
	defer func() {
		log.Printf("closing postgres connection")
		if err := sqlDB.Close(); err != nil {
			log.Printf("postgres close error: %v", err)
		}
	}()
	defer func() {
		log.Printf("closing redis connection")
		if err := redisClient.Close(); err != nil {
			log.Printf("redis close error: %v", err)
		}
	}()

	app := fiber.New(fiber.Config{
		ErrorHandler: httproutes.ErrorHandler(),
	})

	// CORS (Cross-Origin Resource Sharing / Çapraz Kaynak Paylaşımı)
	// Frontend: http://localhost:3000  -> Backend: http://localhost:8080
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://127.0.0.1:3000",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))
	app.Use(fiberlogger.New(fiberlogger.Config{
		Format:     "[${time}] ${status} ${latency} ${method} ${path}\n",
		TimeFormat: time.RFC3339,
		TimeZone:   "Local",
	}))

	httproutes.RegisterAuthRoutes(app, db, cfg.Auth)
	httproutes.RegisterConsentRoutes(app)
	httproutes.RegisterSurveyRoutes(app, db, cfg.Auth.JWTSecret)
	httproutes.RegisterVoteRoutes(app, db, redisClient, cfg.Auth.JWTSecret)
	httproutes.RegisterResultsReportRoutes(app, db, cfg.Auth.JWTSecret)
	registerHealthRoute(
		app,
		func(ctx context.Context) error { return sqlDB.PingContext(ctx) },
		func(ctx context.Context) error { return redisClient.Ping(ctx).Err() },
	)

	address := ":" + cfg.App.Port
	log.Printf("startup complete; API listening on %s", address)

	listenErrCh := make(chan error, 1)
	go func() {
		listenErrCh <- app.Listen(address)
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	select {
	case sig := <-sigCh:
		log.Printf("shutdown signal received: %s", sig)
	case err := <-listenErrCh:
		if err != nil {
			log.Printf("api listen stopped with error: %v", err)
		} else {
			log.Printf("api listener stopped")
		}
		return
	}

	log.Printf("shutting down API server")
	if err := app.Shutdown(); err != nil {
		log.Printf("api shutdown error: %v", err)
	}
	log.Printf("api server stopped")
}
