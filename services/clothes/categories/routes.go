package categories

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/utils"
	"github.com/omkarp02/pro/utils/constant"
	"github.com/omkarp02/pro/utils/validation"
)

type CategoryService interface {
	Create(ctx context.Context, createCategory TCreateCategory) (string, error)
	GetAllCategory(ctx context.Context, filterData FilterCategoryModal, project []string, inclusive bool) ([]TCategoryList, error)
}

type Handler struct {
	service   CategoryService
	cfg       *config.Config
	validator *validation.Validator
}

func NewHandler(service CategoryService, cfg *config.Config, validator *validation.Validator) *Handler {
	return &Handler{service: service, cfg: cfg, validator: validator}
}

func (h *Handler) RegisterRoutes(router router.Router, link string) {
	routeGrp := router.Group(link)

	routeGrp.Post("/", h.create)
	routeGrp.Get("/", h.get)

	routeGrp.Get("/admin", h.getForAdmin)
}

func (h *Handler) create(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	var createCategory TCreateCategory

	if err := h.validator.ValidateBody(c, &createCategory); err != nil {
		return err
	}

	id, err := h.service.Create(ctx, createCategory)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Category Created Successfully", fiber.Map{"id": id}, 201)
}

func (h *Handler) get(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	var filterCat TFilterCategory

	if err := h.validator.ValidateParams(c, &filterCat); err != nil {
		return err
	}
	project := []string{"catId", "icon", "name", "slug"}
	data, err := h.service.GetAllCategory(ctx, FilterCategoryModal{Status: constant.STATUS_ACTIVE, Pagination: filterCat.Pagination}, project, true)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Category Created Successfully", data, 201)
}

func (h *Handler) getForAdmin(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	var filterCat TFilterCategory

	if err := h.validator.ValidateParams(c, &filterCat); err != nil {
		return err
	}

	project := []string{"catId", "name", "slug", "_id"}
	data, err := h.service.GetAllCategory(ctx, FilterCategoryModal{Status: constant.STATUS_ACTIVE, Pagination: filterCat.Pagination}, project, true)
	if err != nil {
		return err
	}
	return utils.SendResponse(c, "Category Created Successfully", data, 201)
}

func createContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
