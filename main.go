package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"

	"sessionReplay/config"
	"sessionReplay/db"
	"sessionReplay/kafka"
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

	//making simple kafka set up for producer and consumer

	// Initialize database
	database, err := db.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer database.Close()

	// Start your app logic here
	log.Println("Database connection established successfully!")
	// app.Get("/api", func(c *fiber.Ctx) error {
	// 	return c.SendString("Hello World")
	// })

	//making simple kafka set up for producer and consumer
	kafka.ProduceMessage("localhost:9092", "test-topic")
	kafka.ConsumeMessages("localhost:9092", "test-topic", "test-group")
	app.Listen(":8080")

}
