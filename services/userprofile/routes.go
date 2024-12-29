package userprofile

import (
	"fmt"
	"log/slog"

	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/services/middleware"
	"github.com/omkarp02/pro/utils"
	"github.com/omkarp02/pro/utils/validation"
)

type UserStore interface {
	CreateUser(CreateUser) (interface{}, error)
	GetUser(string) (*User, error)
}

type Handler struct {
	store     UserStore
	cfg       *config.Config
	validator *validation.Validator
}

func NewHandler(store UserStore, cfg *config.Config, validator *validation.Validator) *Handler {
	return &Handler{store: store, cfg: cfg, validator: validator}
}

func (h *Handler) RegisterRoutes(router router.Router, link string) {
	routeGrp := router.Group(link)
	routeGrp.Use(middleware.VerifyToken(h.cfg))

	routeGrp.Post("/", h.create)
	routeGrp.Get("/", h.get)
}

func (h *Handler) create(c router.Context) error {

	var user CreateUser
	// Parse the JSON body into the struct
	if err := h.validator.ValidateBody(c, &user); err != nil {
		return err
	}

	slog.Debug("created User data routes", "user", user)

	id, err := h.store.CreateUser(user)
	if err != nil {
		return err
	}

	slog.Debug("UserId", "id", fmt.Sprintf("%v", id))
	return utils.SendResponse(c, "User created successfully", nil, 201)
}

func (h *Handler) get(c router.Context) error {

	res, err := h.store.GetUser("skldfjlskdjf")
	if err != nil {
		slog.Error("error", "err", err)
		return err
	}
	return utils.SendResponse(c, "user retrieved", res, 200)
}
