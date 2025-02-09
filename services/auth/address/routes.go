package address

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/services/middleware"
	"github.com/omkarp02/pro/services/utils/store"
	"github.com/omkarp02/pro/utils"
	"github.com/omkarp02/pro/utils/validation"
)

type AddressService interface {
	Create(ctx context.Context, userId string, createPayload TCreateAddress) (string, error)
	GetAddressByUserId(ctx context.Context, userId string) ([]Address, error)
	UpdateAddress(ctx context.Context, payload UpdateAddressModel) error
	DeleteAddress(ctx context.Context, addressId string, userId string) error
	DeleteAddressByIds(ctx context.Context, addressIds []string, userId string) error
	FindById(ctx context.Context, id string, userId string) (Address, error)
	FindPrimaryAddress(ctx context.Context, userId string) (Address, error)
}

type Handler struct {
	service   AddressService
	cfg       *config.Config
	validator *validation.Validator
}

func NewHandler(service AddressService, cfg *config.Config, validator *validation.Validator) *Handler {
	return &Handler{service: service, cfg: cfg, validator: validator}
}

func (h *Handler) RegisterRoutes(router router.Router, link string) {
	routeGrp := router.Group(link)

	routeGrp.Use(middleware.VerifyToken(h.cfg))

	routeGrp.Post("/", h.create)
	routeGrp.Delete("/:id", h.deleteUserAddress)
	routeGrp.Get("/is-primary", h.getPrimaryAddress)
	routeGrp.Get("/:id", h.getAddress)
	routeGrp.Post("/delete/many", h.deleteUserAddressByUserIds)
	routeGrp.Put("/", h.UpdateUserAddress)
	routeGrp.Patch("/is-primary", h.ChangeAddressIsPrimary)
	routeGrp.Get("/", h.getByUserId)
}

func (h *Handler) create(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	var createAddress TCreateAddress

	if err := h.validator.ValidateBody(c, &createAddress); err != nil {
		return err
	}

	userId := c.GetDecodedData().ID

	id, err := h.service.Create(ctx, userId, createAddress)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Address Created Successfully", fiber.Map{"id": id}, 201)
}

func (h *Handler) getByUserId(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	decodedUserId := c.GetDecodedData().ID

	data, err := h.service.GetAddressByUserId(ctx, decodedUserId)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Address Fetched Successfully", data, 200)
}

func (h *Handler) getAddress(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	addressId, err := h.validator.ValidateParam(c, "id")
	if err != nil {
		return err
	}

	decodedUserId := c.GetDecodedData().ID

	data, err := h.service.FindById(ctx, addressId, decodedUserId)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Address Fetched Successfully", data, 200)
}

func (h *Handler) getPrimaryAddress(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()
	fmt.Println(">>>>>>>>>> here rezched")

	decodedUserId := c.GetDecodedData().ID

	data, err := h.service.FindPrimaryAddress(ctx, decodedUserId)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Address Fetched Successfully", data, 200)
}

func (h *Handler) deleteUserAddress(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	addressId, err := h.validator.ValidateParam(c, "id")
	if err != nil {
		return err
	}

	decodedUserId := c.GetDecodedData().ID

	err = h.service.DeleteAddress(ctx, addressId, decodedUserId)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Address Deleted Successfully", nil, 200)
}

func (h *Handler) deleteUserAddressByUserIds(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	var payload TDeleteAddressByIds

	err := h.validator.ValidateBody(c, &payload)
	if err != nil {
		return err
	}

	decodedUserId := c.GetDecodedData().ID

	err = h.service.DeleteAddressByIds(ctx, payload.Ids, decodedUserId)

	return utils.SendResponse(c, "Address Deleted Successfully", nil, 200)
}

func (h *Handler) UpdateUserAddress(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	var payload TUpdateAddress
	decodedUserId := c.GetDecodedData().ID

	if err := h.validator.ValidateBody(c, &payload); err != nil {
		return err
	}

	model := UpdateAddressModel{
		Id:        payload.Id,
		Address:   store.AddressModel(payload.Address),
		IsPrimary: payload.IsPrimary,
		UserID:    decodedUserId,
	}

	err := h.service.UpdateAddress(ctx, model)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Address Updated Successfully", nil, 200)
}

func (h *Handler) ChangeAddressIsPrimary(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	var payload TUpdateAddressIsPrimary
	decodedUserId := c.GetDecodedData().ID

	if err := h.validator.ValidateBody(c, &payload); err != nil {
		return err
	}

	model := UpdateAddressModel{
		Id:        payload.Id,
		IsPrimary: payload.IsPrimary,
		UserID:    decodedUserId,
	}

	err := h.service.UpdateAddress(ctx, model)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Address Updated Successfully", nil, 200)
}

func createContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
