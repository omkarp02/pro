package cart

import (
	"context"
	"time"

	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/services/middleware"
	"github.com/omkarp02/pro/utils"
	"github.com/omkarp02/pro/utils/validation"
)

type CartService interface {
	AddToCard(ctx context.Context, userId string, cartDetails TAddToCart) error
	UpdateQuantityOfItem(ctx context.Context, userId string, payload IUpdateQuantityOfItem) error
	FindOne(ctx context.Context, userId string) (IFindOneRes, error)
}

type Handler struct {
	service   CartService
	cfg       *config.Config
	validator *validation.Validator
}

func NewHandler(service CartService, cfg *config.Config, validator *validation.Validator) *Handler {
	return &Handler{service: service, cfg: cfg, validator: validator}
}

func (h *Handler) RegisterRoutes(router router.Router, link string) {
	routeGrp := router.Group(link)

	routeGrp.Use(middleware.VerifyToken(h.cfg))
	routeGrp.Post("/add", h.addToCard)
	routeGrp.Patch("/item-quantity", h.updateCartItemQuantity)
	routeGrp.Get("/", h.getCart)
}

func (h *Handler) addToCard(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	var cartDetails TAddToCart

	if err := h.validator.ValidateBody(c, &cartDetails); err != nil {
		return err
	}

	userId := c.GetDecodedData().ID

	err := h.service.AddToCard(ctx, userId, cartDetails)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Cart Created Successfully", nil, 200)
}

func (h *Handler) updateCartItemQuantity(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	var cartDetails IUpdateQuantityOfItem

	userId := c.GetDecodedData().ID

	if err := h.validator.ValidateBody(c, &cartDetails); err != nil {
		return err
	}

	err := h.service.UpdateQuantityOfItem(ctx, userId, cartDetails)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Cart Updated Successfully", nil, 200)
}

func (h *Handler) getCart(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	userId := c.GetDecodedData().ID

	cart, err := h.service.FindOne(ctx, userId)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Cart Updated Successfully", cart, 200)
}

func createContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
