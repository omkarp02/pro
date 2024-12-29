package errutil

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/omkarp02/pro/types"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Middleware for centralized error handling
func ErrorHandler(c *fiber.Ctx, err error) error {

	// Default to Internal Server Error
	statusCode := http.StatusInternalServerError
	code := 1
	msg := err.Error()

	if apiErr, ok := err.(APIError); ok {
		statusCode = apiErr.StatusCode
		code = apiErr.Code
		msg = apiErr.Msg

	}

	// Log the error
	slog.Error("Error occurred", "err", err.Error())

	// Send a standardized JSON response
	return c.Status(statusCode).JSON(types.Response{
		Status:  code,
		Message: msg,
	})
}

func handleMongoError(err error) (string, bool) {

	var msg string

	// // Check for MongoDB errors
	// var mongoErr mongo.CommandError
	// if errors.As(err, &mongoErr) {
	// 	code = fiber.StatusBadRequest
	// 	msg = "MongoDB Command Error: " + mongoErr.Message
	// }

	// var writeErr mongo.WriteException
	// if errors.As(err, &writeErr) {
	// 	code = fiber.StatusConflict
	// 	msg = "MongoDB Write Error: " + writeErr.Error()
	// }

	// Check for MongoDB errors
	if mongo.IsDuplicateKeyError(err) {
		msg = "Duplicate key error in MongoDB"
	} else if errors.Is(err, mongo.ErrNoDocuments) {
		msg = "Document not found"
	} else {
		return "", false
	}

	return msg, true
}
