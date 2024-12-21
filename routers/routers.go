package routers

import (
	"github.com/gofiber/fiber/v2"

	"sessionReplay/handlers"
)

func SetupRoutes(app *fiber.App) {
	// API routes
	app.Get("/health", handlers.HealthCheck)
	app.Get("/health", handlers.HealthCheck)

}
