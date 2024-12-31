package product

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/utils"
	"github.com/omkarp02/pro/utils/validation"
)

type ProductService interface {
	CreateProductList(ctx context.Context, createProductList TCreateProductList) (string, error)
	FilterProductList(ctx context.Context, filterProductList TFilterProductList) ([]ProductList, error)
}

type Handler struct {
	service   ProductService
	cfg       *config.Config
	validator *validation.Validator
}

func NewHandler(service ProductService, cfg *config.Config, validator *validation.Validator) *Handler {
	return &Handler{service: service, cfg: cfg, validator: validator}
}

func (h *Handler) RegisterRoutes(router router.Router, link string) {
	routeGrp := router.Group(link)

	routeGrp.Post("/create/product-list", h.createProdutList)
	routeGrp.Get("/filter/product-list", h.getFilteredProductList)
}

func (h *Handler) createProdutList(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	var productList TCreateProductList

	if err := h.validator.ValidateBody(c, &productList); err != nil {
		fmt.Println(err)
		return err
	}

	id, err := h.service.CreateProductList(ctx, productList)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return utils.SendResponse(c, "Product List Created Successfully", fiber.Map{"id": id}, 201)
}

func (h *Handler) getFilteredProductList(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()
	var filterData TFilterProductList

	if err := h.validator.ValidateParams(c, &filterData); err != nil {
		return err
	}

	data, err := h.service.FilterProductList(ctx, filterData)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Product List Created Successfully", data, 200)
}

func createContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
