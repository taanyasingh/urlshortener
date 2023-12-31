package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/tanyasingh/urlshortener/routes"
)

func setRoutes(app *fiber.App) {
	app.Get("/:url", routes.RedirectAPI)
	app.Post("/api/v1", routes.ShortenURL)
	app.Get("/api/v1/metrics", routes.MetricAPI)
}

func main() {
	fmt.Println("Starting server!!!!")
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	app := fiber.New()

	app.Use(logger.New())

	setRoutes(app)

	log.Fatal(app.Listen(os.Getenv("APP_PORT")))

}
