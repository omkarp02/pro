package middleware

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/types"
	"github.com/omkarp02/pro/utils"
)

func VerifyToken(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")
		if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			return utils.UnAuthorized("Authorization header is missing or invalid")
		}

		// Extract the token part
		tokenString := authHeader[7:]

		accessTokenGenerator, err := utils.TokenFactory(types.ACCESS_TOKEN, cfg)
		if err != nil {
			return err
		}

		data, err := accessTokenGenerator.ValidateToken(tokenString)
		if err != nil {
			slog.Error("error while validating token", "error", err)
			return utils.UnAuthorized("Invalid Token")
		}

		userData, err := utils.GetUserDataFromAccessClaimsData(data)
		if err != nil {
			return err
		}

		c.Locals("user", userData)
		return c.Next()

	}
}
