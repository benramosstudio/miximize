package main

import (
	"log"

	"github.com/benramosstudio/miximize/services"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()
	rService := services.NewRService("scripts/robyn.r")

	app.Post("/run-r", func(c fiber.Ctx) error {
		var req services.RRequest
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

		return c.Send(results)
	})

	log.Fatal(app.Listen(":4000"))
}
