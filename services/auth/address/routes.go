package address

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/services/middleware"
	"github.com/omkarp02/pro/utils"
	"github.com/omkarp02/pro/utils/validation"
)

type AddressService interface {
	Create(ctx context.Context, userId string, createPayload TCreateAddress) (string, error)
	GetAddressByUserId(ctx context.Context, userId string) ([]Address, error)
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
	routeGrp.Get("/get-by-userid", h.GetByUserId)
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

func (h *Handler) GetByUserId(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	decodedUserId := c.GetDecodedData().ID

	data, err := h.service.GetAddressByUserId(ctx, decodedUserId)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Address Fetched Successfully", data, 200)
}

func createContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
