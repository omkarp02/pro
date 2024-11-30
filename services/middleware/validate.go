package middleware

// func ValidateBody(objType interface{}) fiber.Handler {
// 	return func(c *fiber.Ctx) error {

// 		obj := reflect.New(reflect.TypeOf(objType)).Interface()

// 		if err := c.BodyParser(&obj); err != nil {
// 			return utils.InvalidReqBody()
// 		}

// 		// Retrieve the validator instance from the context
// 		validate, ok := c.Locals("validator").(*validator.Validate)
// 		if !ok {
// 			return utils.GenerateError(fiber.StatusInternalServerError, "Validator not found")
// 		}

// 		fmt.Println(obj)

// 		if err := validate.Struct(obj); err != nil {
// 			slog.Error("err", "err", err)
// 			return utils.HandleValidationError(err)
// 		}

// 		c.Locals("body", obj)
// 		return c.Next()

// 	}
// }
