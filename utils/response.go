package utils

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/omkarp02/pro/router"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SendResponse(c router.Context, msg string, data interface{}, status int) error {
	//Validating that the status code is in the range 2xx
	if status < 200 || status > 299 {
		slog.Error("Invalid success status code provided", "status", status)
		return c.JSON(fiber.StatusInternalServerError, Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Invalid success status code",
		})
	}

	response := Response{
		Status:  -1,
		Message: msg,
		Data:    data,
	}

	return c.JSON(status, response)
}
