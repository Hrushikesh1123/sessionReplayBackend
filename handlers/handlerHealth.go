package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func HealthCheck(c *fiber.Ctx) error {
	fmt.Println("status check called")
	var a int
	var b int
	fmt.Println(a + b)
	fmt.Println("status check called")

	return c.Status(200).JSON(fiber.Map{
		"status": "ok",
	})
}
