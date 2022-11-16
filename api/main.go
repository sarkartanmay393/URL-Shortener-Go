package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/sarkartanmay393/URL-Shortener-Go/routes"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()
	app.Use(logger.New())

	err = setupRoutes(app)
	if err != nil {
		log.Println(err)
	}

	log.Fatalln(app.Listen(os.Getenv("APP_PORT")))
}

func setupRoutes(app *fiber.App) error {
	app.Get("/:id", routes.ResolveURL)
	app.Post("/api/v1", routes.ShortenURL)

	return nil
}
