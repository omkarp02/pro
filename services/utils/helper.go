package services

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/omkarp02/pro/types"
	"golang.org/x/crypto/bcrypt"
)

var providerIdAndName = make(map[string]string)

func HashPassword(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), 10)
	return string(bytes), err
}

func CheckPasswordHash(pass string, hashPass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(pass))
	return err == nil
}

func UpdateRefreshTokenCookie(c *fiber.Ctx, cookieName string, value string) {
	c.Cookie(&fiber.Cookie{
		Name:    cookieName,
		Value:   value,
		Expires: types.REFRESH_TOKEN_COOKIE_EXPIRY,
		// HTTPOnly: true,
		// Secure:   true,
		// SameSite: fiber.CookieSameSiteStrictMode,
	})
}

func ClearCookie(c *fiber.Ctx, name string) {
	c.Cookie(&fiber.Cookie{
		Name:    name,
		Value:   "",
		Expires: time.Now().Add(-time.Hour),
	})
}

func CreateAccessTokenPayload(id string, providerId string) types.ACCESS_TOKEN_PAYLOAD {
	return types.ACCESS_TOKEN_PAYLOAD{
		ID:         id,
		ProviderId: providerId,
	}
}

func CreateRefreshTokenPayload(id string, providerId string) types.REFRESH_TOKEN_PAYLOAD {
	return types.REFRESH_TOKEN_PAYLOAD{
		ID:         id,
		ProviderId: providerId,
	}
}

func ValidateDataForAccessToken(data interface{}) types.ACCESS_TOKEN_PAYLOAD {
	//here make function to check the data coming from local("user") is valid
	if data, ok := data.(types.ACCESS_TOKEN_PAYLOAD); ok {
		return data
	}

	panic("invalid data")
}
