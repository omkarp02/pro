package state

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/utils"
	"github.com/omkarp02/pro/utils/validation"
)

type StateService interface {
	Create(ctx context.Context, createBody CreateModal) (string, error)
	Find(ctx context.Context, filterPayload FilterListModel) ([]State, error)
	FindById(ctx context.Context, id string) (State, error)
}

type Handler struct {
	service   StateService
	cfg       *config.Config
	validator *validation.Validator
}

func NewHandler(service StateService, cfg *config.Config, validator *validation.Validator) *Handler {
	return &Handler{service: service, cfg: cfg, validator: validator}
}

func (h *Handler) RegisterRoutes(router router.Router, link string) {
	routeGrp := router.Group(link)

	routeGrp.Get("/state", h.get)

	// routeGrp.Use(middleware.VerifyToken(h.cfg))
	// routeGrp.Use(middleware.IsAdmin())

	routeGrp.Post("/state", h.create)
}

func (h *Handler) create(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	var createBody TCreate

	if err := h.validator.ValidateBody(c, &createBody); err != nil {
		return err
	}

	id, err := h.service.Create(ctx, CreateModal(createBody))
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Data Created Successfully", fiber.Map{"id": id}, 201)
}

func (h *Handler) GetByUserId(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	decodedUserId := c.GetDecodedData().ID

	data, err := h.service.FindById(ctx, decodedUserId)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Data Fetched Successfully", data, 200)
}

func (h *Handler) get(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	var filterData TFilterList

	if err := h.validator.ValidateParams(c, &filterData); err != nil {
		return err
	}
	data, err := h.service.Find(ctx, FilterListModel(filterData))
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Data Fetched Successfully", data, 200)
}

func createContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
