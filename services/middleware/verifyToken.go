package middleware

import (
	"log/slog"

	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/utils"
	"github.com/omkarp02/pro/utils/constant"
	"github.com/omkarp02/pro/utils/errutil"
)

func VerifyToken(cfg *config.Config) router.Handler {
	return func(c router.Context) error {

		authHeader := c.Get("Authorization")
		if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			return errutil.UnAuthorized("Authorization header is missing or invalid")
		}

		// Extract the token part
		tokenString := authHeader[7:]

		accessTokenGenerator, err := utils.TokenFactory(constant.ACCESS_TOKEN, cfg)
		if err != nil {
			return err
		}

		data, err := accessTokenGenerator.ValidateToken(tokenString)
		if err != nil {
			slog.Error("error while validating token", "error", err)
			return errutil.UnAuthorized("Invalid Token")
		}

		userData, err := utils.GetUserDataFromAccessClaimsData(data)
		if err != nil {
			return err
		}

		c.Locals("user", userData)
		return c.Next()
	}
}
