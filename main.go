package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

// Checks if envs are set in os or .env file
func getEnvs() (string, string) {
	host := os.Getenv("SERVER_HOST")
	port := os.Getenv("SERVER_PORT")

	if host == "" && port == "" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Environment variables not set")
		}

		host = os.Getenv("SERVER_HOST")
		port = os.Getenv("SERVER_PORT")
	}

	return host, port
}


func main() {
	// Set environmental variables
	host, port := getEnvs()

	// TODO Add db

	app := fiber.New()

  app.Get("/", func(c *fiber.Ctx) error {
    return c.SendString("Hello, World!")
  })

	// Run Server
	log.Fatal(app.Listen(fmt.Sprintf("%s:%s", host, port)))
}
