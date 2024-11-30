package middleware

import (
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/omkarp02/pro/config"
	services "github.com/omkarp02/pro/services/utils"
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

		claimsMap, ok := data.(map[string]interface{})
		if !ok {
			return utils.UnAuthorized("Invalid token claims format")
		}

		asdf := services.ACCESS_TOKEN_PAYLOAD{
			ID: claimsMap["ID"].(string),
		}

		fmt.Println(asdf)

		c.Locals("user", "lskdfjlsdkfj")
		return c.Next()

	}
}
