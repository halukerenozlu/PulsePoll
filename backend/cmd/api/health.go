package main

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

func registerHealthRoute(
	app *fiber.App,
	pingPostgres func(context.Context) error,
	pingRedis func(context.Context) error,
) {
	app.Get("/health", func(c *fiber.Ctx) error {
		healthCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		postgresOK := pingPostgres(healthCtx) == nil
		redisOK := pingRedis(healthCtx) == nil

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
}
