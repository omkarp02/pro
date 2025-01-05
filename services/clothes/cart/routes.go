package cart

import (
	"context"
	"time"

	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/utils"
	"github.com/omkarp02/pro/utils/validation"
)

type CartService interface {
	AddToCard(ctx context.Context, cartDetails TAddToCart) error
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

	routeGrp.Post("/add-to-cart", h.addToCard)
}

func (h *Handler) addToCard(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	var cartDetails TAddToCart

	cartDetails.UserId = "67776e638a1cf8d369ebc97e"

	if err := h.validator.ValidateBody(c, &cartDetails); err != nil {
		return err
	}

	err := h.service.AddToCard(ctx, cartDetails)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Cart Created Successfully", nil, 200)
}

func createContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
