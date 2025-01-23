package bussiness

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/services/middleware"
	"github.com/omkarp02/pro/services/utils/store"
	"github.com/omkarp02/pro/utils"
	"github.com/omkarp02/pro/utils/validation"
)

type BussinessService interface {
	Create(ctx context.Context, createBody CreateModal, userId string) (string, error)
	Find(ctx context.Context, filterPayload FilterListModel) ([]Business, error)
	FindById(ctx context.Context, id string) (Business, error)
}

type Handler struct {
	service   BussinessService
	cfg       *config.Config
	validator *validation.Validator
}

func NewHandler(service BussinessService, cfg *config.Config, validator *validation.Validator) *Handler {
	return &Handler{service: service, cfg: cfg, validator: validator}
}

func (h *Handler) RegisterRoutes(router router.Router, link string) {
	routeGrp := router.Group(link)

	routeGrp.Use(middleware.VerifyToken(h.cfg))
	routeGrp.Use(middleware.IsOwner())

	routeGrp.Post("/in-owner", h.addBussinessToOwner)
}

func (h *Handler) addBussinessToOwner(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	var createBody TCreate

	if err := h.validator.ValidateBody(c, &createBody); err != nil {
		return err
	}

	userId := c.GetDecodedData().ID

	formattedModal := CreateModal{
		Name:        createBody.Name,
		OwnerID:     userId,
		Category:    createBody.Category,
		Description: createBody.Description,
		Address:     store.AddressModel(createBody.Address),
		Contacts:    createBody.Contacts,
		Website:     createBody.Website,
		LogoUrl:     createBody.LogoUrl,
		CreatorId:   userId,
	}

	id, err := h.service.Create(ctx, formattedModal, userId)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Address Created Successfully", fiber.Map{"id": id}, 201)
}

func (h *Handler) GetByUserId(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	decodedUserId := c.GetDecodedData().ID

	data, err := h.service.FindById(ctx, decodedUserId)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Address Fetched Successfully", data, 200)
}

func (h *Handler) Find(c router.Context) error {
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

	return utils.SendResponse(c, "Address Fetched Successfully", data, 200)
}

func createContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
