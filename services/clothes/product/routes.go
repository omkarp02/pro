package product

import (
	"context"
	"time"

	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/services/middleware"
	"github.com/omkarp02/pro/utils"
	"github.com/omkarp02/pro/utils/validation"
)

type ProductService interface {
	// CreateProductList(ctx context.Context, createProductList TCreateProductList) (string, error)
	FilterProductList(ctx context.Context, filterProductList TFilterProductList) ([]TFilteredProductList, error)
	AddProductsToCollection(ctx context.Context, productData TAddProductToCollection) error
	CreateProduct(ctx context.Context, productDetails TCreateProduct, creatorId string) error
	GetProductDetails(ctx context.Context, id string) (ProductDetail, error)
	CreateProductBatch(ctx context.Context, createPayload TCreateProductBatch, userId string) (string, error)
	FindProductBatch(ctx context.Context, filterPayload FilterProductBatchListModel) ([]ProductBatch, error)
	GetProductBatchDetails(ctx context.Context, code string) (ProductBatch, error)
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

	// routeGrp.Post("/create/product-list", h.createProductList)

	routeGrp.Get("/list", h.getFilteredProductList)
	routeGrp.Get("/details/:productId", h.getProductDetails)
	routeGrp.Get("/batch/:batchId", h.getProductBatchDetails)

	routeGrp.Use(middleware.VerifyToken(h.cfg))
	routeGrp.Use(middleware.IsOwner())

	routeGrp.Get("/batch", h.findProductBatch)
	routeGrp.Post("/", h.createProduct)
	routeGrp.Post("/batch", h.createProductBatch)
	routeGrp.Get("/add-to-collection", h.addToCollection)
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

func (h *Handler) createProduct(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()
	var productData TCreateProduct

	userId := c.GetDecodedData().ID

	if err := h.validator.ValidateBody(c, &productData); err != nil {
		return err
	}

	err := h.service.CreateProduct(ctx, productData, userId)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Product List Created Successfully", "", 200)
}

func (h *Handler) createProductBatch(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	var payload TCreateProductBatch

	userId := c.GetDecodedData().ID

	if err := h.validator.ValidateBody(c, &payload); err != nil {
		return err
	}

	id, err := h.service.CreateProductBatch(ctx, payload, userId)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Product Batch Created Successfully", id, 200)
}

func (h *Handler) findProductBatch(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	var filterData TFilterProductBatchList

	userId := c.GetDecodedData().ID

	if err := h.validator.ValidateParams(c, &filterData); err != nil {
		return err
	}

	payload := FilterProductBatchListModel{
		CreatorId: userId,
		Page:      filterData.Page,
		Limit:     filterData.Limit,
	}

	data, err := h.service.FindProductBatch(ctx, payload)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Data Fetched Successfully", data, 200)
}

func (h *Handler) getProductDetails(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	productId := c.Params("productId")

	productDetails, err := h.service.GetProductDetails(ctx, productId)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Product Detail fetched Successfully", productDetails, 200)
}

func (h *Handler) getProductBatchDetails(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	batchId := c.Params("batchId")

	batchDetails, err := h.service.GetProductBatchDetails(ctx, batchId)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Batch Detail fetched Successfully", batchDetails, 200)
}

func createContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func (h *Handler) addToCollection(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()
	var productData TAddProductToCollection

	if err := h.validator.ValidateBody(c, &productData); err != nil {
		return err
	}

	err := h.service.AddProductsToCollection(ctx, productData)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Product List Created Successfully", "", 200)
}

// func (h *Handler) createProductList(c router.Context) error {
// 	ctx, cancel := createContext()
// 	defer cancel()

// 	var productList TCreateProductList

// 	if err := h.validator.ValidateBody(c, &productList); err != nil {
// 		fmt.Println(err)
// 		return err
// 	}

// 	id, err := h.service.CreateProductList(ctx, productList)
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}

// 	return utils.SendResponse(c, "Product List Created Successfully", fiber.Map{"id": id}, 201)
// }
