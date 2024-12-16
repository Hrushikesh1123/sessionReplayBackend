package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"

	"sessionReplay/config"
	"sessionReplay/db"
	"sessionReplay/routers"
)

func main() {
	fmt.Println("Hello World")

	app := fiber.New()
	routers.SetupRoutes(app)
	fmt.Println("app started", app.Config().StreamRequestBody)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize database
	database, err := db.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer database.Close()

	// Start your app logic here
	log.Println("Database connection established successfully!")
	app.Listen(":8080")
	// app.Get("/api", func(c *fiber.Ctx) error {
	// 	return c.SendString("Hello World")
	// })

}
