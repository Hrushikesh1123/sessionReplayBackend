package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"sessionReplay/routers"
)

func main() {
	fmt.Println("Hello World")

	app := fiber.New()
	routers.SetupRoutes(app)
	fmt.Println("app started", app.Config().StreamRequestBody)

	app.Listen(":8080")
	// app.Get("/api", func(c *fiber.Ctx) error {
	// 	return c.SendString("Hello World")
	// })

}
