package userprofile

import (
	"context"
	"time"

	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/services/middleware"
	"github.com/omkarp02/pro/utils"
	"github.com/omkarp02/pro/utils/validation"
)

type UserService interface {
	CreateUser(ctx context.Context, createUserPayload TCreateUser, userAccountId string) (string, error)
}

type Handler struct {
	service   UserService
	cfg       *config.Config
	validator *validation.Validator
}

func NewHandler(service UserService, cfg *config.Config, validator *validation.Validator) *Handler {
	return &Handler{service: service, cfg: cfg, validator: validator}
}

func (h *Handler) RegisterRoutes(router router.Router, link string) {
	routeGrp := router.Group(link)
	routeGrp.Use(middleware.VerifyToken(h.cfg))

	routeGrp.Post("/", h.create)
	// routeGrp.Get("/", h.get)
}

func (h *Handler) create(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	decodedUserId := c.GetDecodedData().ID

	var user TCreateUser

	// Parse the JSON body into the struct
	if err := h.validator.ValidateBody(c, &user); err != nil {
		return err
	}

	id, err := h.service.CreateUser(ctx, user, decodedUserId)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "User created successfully", id, 201)
}

// func (h *Handler) get(c router.Context) error {
// 	ctx, cancel := createContext()
// 	defer cancel()

// 	res, err := h.service.GetUser("skldfjlskdjf")
// 	if err != nil {
// 		slog.Error("error", "err", err)
// 		return err
// 	}
// 	return utils.SendResponse(c, "user retrieved", res, 200)
// }

func createContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
