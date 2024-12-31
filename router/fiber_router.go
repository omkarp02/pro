package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/omkarp02/pro/types"
)

type FiberRouter struct {
	router fiber.Router
}

func NewFiberRouter(router fiber.Router) *FiberRouter {
	return &FiberRouter{router: router}
}

func (r *FiberRouter) Group(path string) Router {
	return &FiberRouter{router: r.router.Group(path)}
}

func (r *FiberRouter) Post(path string, handler func(ctx Context) error) {
	r.router.Post(path, func(c *fiber.Ctx) error {
		return handler(&FiberContext{c})
	})
}

func (r *FiberRouter) Get(path string, handler func(ctx Context) error) {
	r.router.Get(path, func(c *fiber.Ctx) error {
		return handler(&FiberContext{c})
	})
}

func (r *FiberRouter) Use(middleware ...func(ctx Context) error) {
	for _, m := range middleware {
		r.router.Use(func(c *fiber.Ctx) error {
			return m(&FiberContext{c})
		})
	}
}

// FiberContext implements the Context interface for Fiber
type FiberContext struct {
	ctx *fiber.Ctx
}

func (c *FiberContext) Bind(v interface{}) error {
	return c.ctx.BodyParser(v)
}

func (c *FiberContext) JSON(statusCode int, data interface{}) error {
	return c.ctx.Status(statusCode).JSON(data)
}

func (c *FiberContext) GetDecodedData() types.ACCESS_TOKEN_PAYLOAD {

	userDetails := c.ctx.Locals("users")

	//here make function to check the data coming from local("user") is valid
	if data, ok := userDetails.(types.ACCESS_TOKEN_PAYLOAD); ok {
		return data
	}

	panic("invalid data")
}

func (c *FiberContext) Get(key string) string {
	return c.ctx.Get(key)
}

func (c *FiberContext) Locals(key interface{}, value ...interface{}) {
	c.ctx.Locals(key, value)
}

func (c *FiberContext) Next() error {
	return c.ctx.Next()
}

func (c *FiberContext) GetCookie(name string) string {
	return c.ctx.Cookies(name)
}

func (c *FiberContext) SetCookie(cookie *fiber.Cookie) {
	c.ctx.Cookie(cookie)
}

func (c *FiberContext) GetContext() *fiber.Ctx {
	return c.ctx
}

func (c *FiberContext) Params(key string) string {
	return c.ctx.Params(key)
}

func (c *FiberContext) Redirect(location string, status ...int) error {
	return c.ctx.Redirect(location, status...)
}

func (c *FiberContext) QueryParser(out interface{}) error {

	if err := c.ctx.QueryParser(out); err != nil {
		return err
	}
	return nil

}
