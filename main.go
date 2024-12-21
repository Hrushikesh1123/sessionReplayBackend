package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"

	"sessionReplay/config"
	"sessionReplay/db"
	"sessionReplay/kafka"
	"sessionReplay/routers"
)

func main() {
	fmt.Println("Hello World")

	// 1. Create the Fiber instance (ONE TIME)
	app := fiber.New()

	// 2. Set up routes on the Fiber instance
	routers.SetupRoutes(app)
	fmt.Println("app started", app.Config().StreamRequestBody)

	// 3. Load configuration (DB, Kafka, etc.)
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// 4. Initialize database
	database, err := db.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer database.Close()
	log.Println("Database connection established successfully!")

	// 5. Initialize Kafka client
	kafkaClient := kafka.NewClient(kafka.KafkaInitParams{
		Brokers: cfg.KafkaBrokers,  // e.g. ["localhost:9092"]
		Topics:  cfg.KafkaTopics,   // e.g. ["topicA","topicB"]
		Groups:  cfg.KafkaGroupIDs, // e.g. ["groupA","groupB"]
	})
	defer kafkaClient.Close()

	// 6. Start Kafka consumers in the background
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	kafkaClient.StartAllConsumers(ctx, &wg)

	// 7. Example route for producing Kafka messages
	app.Post("/produce", func(c *fiber.Ctx) error {
		type payload struct {
			Topic string `json:"topic"`
			Key   string `json:"key"`
			Value string `json:"value"`
		}
		var p payload
		if err := c.BodyParser(&p); err != nil {
			return err
		}
		if err := kafkaClient.ProduceMessage(p.Topic, p.Key, p.Value); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.SendString("Message produced successfully!")
	})

	// 8. Start Fiber in a separate goroutine (non-blocking)
	go func() {
		if err := app.Listen(":" + "8080"); err != nil {
			log.Fatalf("Fiber failed: %v", err)
		}
	}()

	// 9. Listen for system signals (Ctrl+C, SIGTERM)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("[main] Shutting down...")

	// 10. Cancel the consumer context to stop Kafka consumers gracefully
	cancel()
	wg.Wait()

	// 11. Gracefully shut down Fiber (with a timeout)
	_, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := app.Shutdown(); err != nil {
		log.Printf("[main] Fiber shutdown error: %v", err)
	}

	log.Println("[main] All done!")
}
