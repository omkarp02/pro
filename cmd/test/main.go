package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/services/auth/useraccount"
	"github.com/omkarp02/pro/utils/validation"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type FiberRouter struct {
	router fiber.Router
	app    *fiber.App
	cfg    *config.Config
}

func main() {

	cfg := config.MustLoad("")
	config.SetUpLogger()

	DB, err := db.NewDatabase(cfg)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err = DB.DB.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	api := router.NewFiberRouter(cfg)
	api.Get("/", func(ctx router.Context) error {
		response := Response{
			Status:  -1,
			Message: "User registered successfully",
			Data:    fiber.Map{"id": "sdflgkj"},
		}

		return ctx.JSON(200, response)
	})

	api.Post("/", func(ctx router.Context) error {
		response := Response{
			Status:  -1,
			Message: "User registered successfully",
			Data:    fiber.Map{"id": "sdflgkj"},
		}

		return ctx.JSON(200, response)
	})

	validator := validation.NewValidator()

	useraccount.Intialize(DB, cfg, validator, api)

	api.Listen(cfg.HTTPServer.Addr)
}
