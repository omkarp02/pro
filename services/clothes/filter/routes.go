package filter

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/utils"
	"github.com/omkarp02/pro/utils/validation"
)

type FilterService interface {
	CreateFilter(ctx context.Context, createFilter TCreateFilter) (string, error)
	CreateFilterType(ctx context.Context, createFilterType TCreateFilterType) (string, error)
}

type Handler struct {
	service   FilterService
	cfg       *config.Config
	validator *validation.Validator
}

func NewHandler(service FilterService, cfg *config.Config, validator *validation.Validator) *Handler {
	return &Handler{service: service, cfg: cfg, validator: validator}
}

func (h *Handler) RegisterRoutes(router router.Router, link string) {
	routeGrp := router.Group(link)

	routeGrp.Post("/", h.createFilter)
	routeGrp.Post("/type", h.createFilterType)
}

func (h *Handler) createFilter(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	var filterDetails TCreateFilter

	if err := h.validator.ValidateBody(c, &filterDetails); err != nil {
		return err
	}

	id, err := h.service.CreateFilter(ctx, filterDetails)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Filter Created Successfully", fiber.Map{"id": id}, 201)
}

func (h *Handler) createFilterType(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	var filterTypeDetails TCreateFilterType

	if err := h.validator.ValidateBody(c, &filterTypeDetails); err != nil {
		return err
	}

	id, err := h.service.CreateFilterType(ctx, filterTypeDetails)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Filter Type Created Successfully", fiber.Map{"id": id}, 201)
}

func createContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
