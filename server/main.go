package main

import (
	"log"
	"time"

	"github.com/benramosstudio/miximize/internal/models"
	"github.com/benramosstudio/miximize/services"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New(fiber.Config{
		ReadTimeout:  300 * time.Second, // 5 minutes
		WriteTimeout: 300 * time.Second, // 5 minutes
	})

	rService := services.NewRService("scripts/robyn.r")

	app.Post("/run-r", func(c fiber.Ctx) error {
		var req models.RobynParams
		if err := c.Bind().JSON(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid request body",
			})
		}

		results, err := rService.ProcessRScript(req)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		c.Set("Content-Type", "application/json")
		return c.Send(results)
	})

	log.Fatal(app.Listen(":4000"))
}
