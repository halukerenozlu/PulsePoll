package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	// CORS (Cross-Origin Resource Sharing / Çapraz Kaynak Paylaşımı)
	// Frontend: http://localhost:3000  -> Backend: http://localhost:8080
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000,http://127.0.0.1:3000",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"ok": true})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("API listening on :%s", port)
	log.Fatal(app.Listen(":" + port))
}