package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/omkarp02/pro/types"
)

type Router interface {
	Group(path string) Router
	Post(path string, handler func(ctx Context) error)
	Get(path string, handler func(ctx Context) error)
	Use(middleware ...func(ctx Context) error)
}

type Handler func(ctx Context) error

type Context interface {
	Bind(interface{}) error
	JSON(statusCode int, data interface{}) error
	GetDecodedData() types.ACCESS_TOKEN_PAYLOAD
	Get(key string) string
	Locals(key interface{}, value ...interface{})
	Next() error
	GetCookie(name string) string
	SetCookie(cookie *fiber.Cookie)
	GetContext() *fiber.Ctx
	Params(key string) string
	Redirect(location string, status ...int) error
}
