package api

type Router interface {
	Group(path string) Router
	Post(path string, handler func(ctx Context) error)
	Get(path string, handler func(ctx Context) error)
	Use(middleware ...func(ctx Context) error)
}

type Context interface {
	Bind(interface{}) error
	JSON(statusCode int, data interface{}) error
}
