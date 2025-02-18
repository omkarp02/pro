package router

import (
	"fmt"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/types"
	"github.com/omkarp02/pro/utils/errutil"
)

type FiberRouter struct {
	router fiber.Router
	app    *fiber.App
	cfg    *config.Config
}

func NewFiberRouter(cfg *config.Config) *FiberRouter {

	fiberConfig := fiber.Config{
		ErrorHandler:             errutil.ErrorHandler,
		EnableSplittingOnParsers: true,
	}
	app := fiber.New(fiberConfig)

	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Cors.AllowOrigins,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		ExposeHeaders:    "Content-Length,Content-Type",
		AllowCredentials: true,
		MaxAge:           86400,
	}))
	app.Use(healthcheck.New())

	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: cfg.Secret.CookieEncryptionKey,
	}))

	api := app.Group("/api/v1")

	api.Get("/metrics", monitor.New())
	api.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	return &FiberRouter{router: api, app: app, cfg: cfg}
}

func (r *FiberRouter) Listen(path string) error {
	return r.app.Listen(r.cfg.HTTPServer.Addr)
}

func (r *FiberRouter) Group(path string) Router {
	return &FiberRouter{router: r.router.Group(path)}
}

func (r *FiberRouter) Post(path string, handler func(ctx Context) error) {
	r.router.Post(path, func(c *fiber.Ctx) error {
		return handler(&FiberContext{c})
	})
}

func (r *FiberRouter) Put(path string, handler func(ctx Context) error) {
	r.router.Put(path, func(c *fiber.Ctx) error {
		return handler(&FiberContext{c})
	})
}

func (r *FiberRouter) Delete(path string, handler func(ctx Context) error) {
	r.router.Delete(path, func(c *fiber.Ctx) error {
		return handler(&FiberContext{c})
	})
}

func (r *FiberRouter) Patch(path string, handler func(ctx Context) error) {
	r.router.Patch(path, func(c *fiber.Ctx) error {
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

	userDetails := c.ctx.Locals("user")

	fmt.Println(reflect.TypeOf(userDetails))

	validArray, ok := userDetails.([]interface{})
	if ok {
		if data, ok := validArray[0].(types.ACCESS_TOKEN_PAYLOAD); ok {
			return data
		}
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

func (c *FiberContext) Params(key string, defaultValue ...string) string {
	return c.ctx.Params(key, defaultValue...)
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
