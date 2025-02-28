package helper

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/omkarp02/pro/router"
)

func UpdateCookie(c router.Context, cookieName string, value string, expiry time.Time) {
	c.SetCookie(&fiber.Cookie{
		Name:     cookieName,
		Value:    value,
		Expires:  expiry,
		Path:     "/",
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteNoneMode,
		Domain:   "localhost:3000",
	})
}

func ClearCookie(c router.Context, name string) {
	c.SetCookie(&fiber.Cookie{
		Name:    name,
		Value:   "",
		Expires: time.Now().Add(-time.Hour),
	})
}
