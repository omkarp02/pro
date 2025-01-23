package owner

import (
	"context"
	"time"

	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/services/middleware"
	"github.com/omkarp02/pro/utils"
	"github.com/omkarp02/pro/utils/validation"
)

type TestService interface {
	Create(ctx context.Context, createBody CreateModal) (string, error)
	Find(ctx context.Context, filterPayload FilterOwnerListModel) ([]Owner, error)
}

type Handler struct {
	service   TestService
	cfg       *config.Config
	validator *validation.Validator
}

func NewHandler(service TestService, cfg *config.Config, validator *validation.Validator) *Handler {
	return &Handler{service: service, cfg: cfg, validator: validator}
}

func (h *Handler) RegisterRoutes(router router.Router, link string) {
	routeGrp := router.Group(link)

	routeGrp.Use(middleware.VerifyToken(h.cfg))
	routeGrp.Use(middleware.IsAdmin())

	routeGrp.Get("/", h.GetAll)
}

func (h *Handler) GetAll(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	var filterData TFilterList

	if err := h.validator.ValidateParams(c, &filterData); err != nil {
		return err
	}
	data, err := h.service.Find(ctx, FilterOwnerListModel(filterData))
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Address Fetched Successfully", data, 200)
}

// func (h *Handler) create(c router.Context) error {
// 	ctx, cancel := createContext()
// 	defer cancel()

// 	var createBody TCreateOwner

// 	if err := h.validator.ValidateBody(c, &createBody); err != nil {
// 		return err
// 	}

// 	userId := c.GetDecodedData().ID

// 	modal := CreateModal{
// 		TCreateOwner: createBody,
// 		CreatorId:    userId,
// 	}

// 	id, err := h.service.Create(ctx, modal)
// 	if err != nil {
// 		return err
// 	}

// 	return utils.SendResponse(c, "Address Created Successfully", fiber.Map{"id": id}, 201)
// }

func createContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
