package filter

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/services/middleware"
	"github.com/omkarp02/pro/utils"
	"github.com/omkarp02/pro/utils/constant"
	"github.com/omkarp02/pro/utils/validation"
)

type FilterService interface {
	CreateFilterType(ctx context.Context, createFilterType CreateFilterTypeModal) (string, error)
	CreateFilter(ctx context.Context, createFilter CreateFilterModal) (string, error)
	FindFitlerType(ctx context.Context, filterPayload FilterTypeListModel) ([]FilterType, error)
	FindFitler(ctx context.Context, filterPayload FilterListModel) ([]Filter, error)
	FindFitlerForUser(ctx context.Context, filterPayload FilterListModel) ([]FilterListForUserRes, error)
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

	routeGrp.Get("/type", h.findFitlerType)
	routeGrp.Get("/", h.findFitler)

	routeGrp.Use(middleware.VerifyToken(h.cfg))
	routeGrp.Use(middleware.IsAdmin())

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

	id, err := h.service.CreateFilter(ctx, CreateFilterModal{
		Name:     filterDetails.Name,
		Type:     filterDetails.Type,
		Category: filterDetails.Category,
		Status:   constant.STATUS_ACTIVE,
		Slug:     filterDetails.Slug,
	})
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

	modal := CreateFilterTypeModal{Name: filterTypeDetails.Name, Status: constant.STATUS_ACTIVE}

	id, err := h.service.CreateFilterType(ctx, modal)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Filter Type Created Successfully", fiber.Map{"id": id}, 201)
}

func (h *Handler) findFitlerType(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	var filterData TFilterTypeList

	if err := h.validator.ValidateParams(c, &filterData); err != nil {
		return err
	}

	data, err := h.service.FindFitlerType(ctx, FilterTypeListModel(filterData))
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Address Fetched Successfully", data, 200)
}

func (h *Handler) findFitler(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	var filterData TFilterList

	if err := h.validator.ValidateParams(c, &filterData); err != nil {
		return err
	}
	data, err := h.service.FindFitlerForUser(ctx, FilterListModel(filterData))
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Address Fetched Successfully", data, 200)
}

func createContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
