package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func main() {
	app := fiber.New()
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Post("/api/v1", func(c *fiber.Ctx) error {
		response := Response{
			Status:  -1,
			Message: "User registered successfully",
			Data:    fiber.Map{"id": "sdflgkj"},
		}

		return c.Status(201).JSON(response)
	})

	app.Listen("0.0.0.0:8080")
}
