package utils

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Middleware for centralized error handling
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Default to Internal Server Error
	code := fiber.StatusInternalServerError
	msg := "Internal Server Error"

	// Check if it's a Fiber error
	if fiberErr, ok := err.(*fiber.Error); ok {
		code = fiberErr.Code
		msg = fiberErr.Message
	}

	if mongoErrMsg, ok := handleMongoError(err); ok {
		code = fiber.StatusInternalServerError
		msg = mongoErrMsg
	}

	// Log the error
	slog.Error("Error occurred", "status", code, "message", msg, "err", err.Error())

	// Send a standardized JSON response
	return c.Status(code).JSON(Response{
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

func GenerateError(code int, msg string) error {
	return fiber.NewError(code, msg)
}

func InvalidCredentails() error {
	return fiber.NewError(fiber.StatusUnauthorized, "Invalid Credentials")
}

func InternalServerError() error {
	return fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")
}

func InvalidReqBody() error {
	return fiber.NewError(fiber.StatusUnprocessableEntity, "Invalid Request Body")
}

func HandleValidationError(errs error) error {

	var errMsgs []string

	for _, err := range errs.(validator.ValidationErrors) {
		field := err.Field()

		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required field", field))
		case "email":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s should be a valid email", field))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is invalid", field))
		}
	}

	return fiber.NewError(fiber.StatusUnprocessableEntity, strings.Join(errMsgs, ","))
}
