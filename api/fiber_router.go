package api

import "github.com/gofiber/fiber/v2"

type FiberRouter struct {
	router fiber.Router
}

func NewFiberRouter(router fiber.Router) *FiberRouter {
	return &FiberRouter{router: router}
}

func (r *FiberRouter) Group(path string) *FiberRouter {
	return NewFiberRouter(r.router.Group(path))
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
