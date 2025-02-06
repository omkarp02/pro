package cart

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/services/middleware"
	"github.com/omkarp02/pro/utils"
	"github.com/omkarp02/pro/utils/errutil"
	"github.com/omkarp02/pro/utils/validation"
)

type CartService interface {
	AddToCard(ctx context.Context, userId string, cartDetails TAddToCart) error
	UpdateQuantityOfItem(ctx context.Context, userId string, payload IUpdateQuantityOfItem) error
	FindOne(ctx context.Context, userId string) (IFindOneRes, error)
	GetTotalItems(ctx context.Context, userId string) (int, error)
	GetCartItemForOffline(ctx context.Context, cartDetails GetCartOfflineModal) ([]IFindOneResCartItem, error)
	DeleteCartItem(ctx context.Context, productCode string, userId string) error
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

	routeGrp.Get("/offline", h.getCartOffline)

	routeGrp.Use(middleware.VerifyToken(h.cfg))
	routeGrp.Post("/add", h.addToCard)
	routeGrp.Patch("/item", h.updateCartItemQuantity)
	routeGrp.Get("/", h.getCart)
	routeGrp.Delete("/item/:cartId", h.deleteCartItem)
	routeGrp.Get("/item/total", h.getCartTotalItem)
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

func (h *Handler) getCartOffline(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	var payload IGetCartOffline

	if err := h.validator.ValidateParams(c, &payload); err != nil {
		return err
	}

	cart, err := h.service.GetCartItemForOffline(ctx, GetCartOfflineModal(payload))
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Cart Updated Successfully", cart, 200)
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

func (h *Handler) deleteCartItem(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	userId := c.GetDecodedData().ID

	code := c.Params("cartId")
	if len(code) == 0 {
		return errutil.InvalidReqData()
	}

	err := h.service.DeleteCartItem(ctx, code, userId)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Cart Updated Successfully", nil, 200)
}

func (h *Handler) getCartTotalItem(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	userId := c.GetDecodedData().ID

	totalItem, err := h.service.GetTotalItems(ctx, userId)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Cart Updated Successfully", fiber.Map{"totalItems": totalItem}, 200)
}

func createContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
