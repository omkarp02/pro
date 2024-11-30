package services

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/omkarp02/pro/utils"
)

func ValidateBody(c *fiber.Ctx, objType interface{}) error {
	obj := objType

	if err := c.BodyParser(obj); err != nil {
		return utils.InvalidReqBody()
	}

	// Retrieve the validator instance from the context
	validate, ok := c.Locals("validator").(*validator.Validate)
	if !ok {
		return utils.GenerateError(fiber.StatusInternalServerError, "Validator not found")
	}

	if err := validate.Struct(obj); err != nil {
		return utils.HandleValidationError(err)
	}

	return nil

}
