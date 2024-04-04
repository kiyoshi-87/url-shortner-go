package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	routes "github.com/kiyoshi-87/url-shortner-API/Routes"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/:url", routes.ResolveURL)
	app.Post("/api/v1", routes.ShortenURL)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("API RUNNING SUCCESSFULLY!")
	})
}

func main() {
	app := fiber.New()
	fmt.Println("Welcome to this url_shortner!")
	app.Use(logger.New()) //USED TO LOG STUFFS

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	SetupRoutes(app)

	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}
