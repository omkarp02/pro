package userprofile

import (
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/utils"
)

type UserStore interface {
	CreateUser(CreateUser) (interface{}, error)
	GetUser(string) (*User, error)
}

type Handler struct {
	store UserStore
	cfg   *config.Config
}

func NewHandler(store UserStore, cfg *config.Config) *Handler {
	return &Handler{store: store, cfg: cfg}
}

func (h *Handler) RegisterRoutes(router fiber.Router, link string) {
	routeGrp := router.Group(link)

	routeGrp.Post("/", h.create)
	routeGrp.Get("/", h.get)
}

func (h *Handler) create(c *fiber.Ctx) error {

	var user CreateUser
	// Parse the JSON body into the struct
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	slog.Debug("created User data routes", "user", user)

	id, err := h.store.CreateUser(user)
	if err != nil {
		return err
	}

	slog.Debug("UserId", "id", fmt.Sprintf("%v", id))
	return utils.SendResponse(c, "User created successfully", nil, 201)
}

func (h *Handler) get(c *fiber.Ctx) error {
	res, err := h.store.GetUser("6747fc459bb5530f8ed389d4")
	if err != nil {
		slog.Error("error", "err", err)
		return err
	}
	return utils.SendResponse(c, "user retrieved", res, 200)
}
