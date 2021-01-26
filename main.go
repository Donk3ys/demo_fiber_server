package main

import (
	handler "fiber_demo/handlers"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)


func main() {
	// Set environmental variables
	host, port := getEnvs()

	// TODO Add db

	// Create App
	app := fiber.New()

	// Route Handler
	gh := handler.GeneralHandler {}
  app.Get("/", gh.SayHi)
	app.Get("/:id", gh.GetPersonMatchingId)
	app.Post("/", gh.PersonCreds)

	// Run Server
	log.Fatal(app.Listen(fmt.Sprintf("%s:%s", host, port)))
}


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
