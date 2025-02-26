package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/omkarp02/pro/utils/errutil"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func main() {

	fiberConfig := fiber.Config{
		ErrorHandler:             errutil.ErrorHandler,
		EnableSplittingOnParsers: true,
	}

	app := fiber.New(fiberConfig)
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://wyse-shop.vercel.app, http://localhost:3000",
		AllowCredentials: true,
	}))

	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: "fe8d78c1e948d78f4d5ef256e4d08c57",
	}))

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
