package product

import (
	"context"
	"strconv"

	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/utils"
	"github.com/omkarp02/pro/utils/constant"
)

type Service struct {
	productTemplateRepo *ProductTemplateRepo
	productListRepo     *ProductListRepo
	productDetailRepo   *ProductDetailRepo
	productBatchRepo    *ProductBatchRepo
	productStatusRepo   *ProductStatusRepo
	txn                 db.TransactionManager
}

func NewService(productListRepo *ProductListRepo, productDetailRepo *ProductDetailRepo, productBatchRepo *ProductBatchRepo, productTemplateRepo *ProductTemplateRepo, txn db.TransactionManager) *Service {
	return &Service{
		productListRepo:     productListRepo,
		productDetailRepo:   productDetailRepo,
		productBatchRepo:    productBatchRepo,
		productTemplateRepo: productTemplateRepo,
		txn:                 txn,
	}
}

func (s *Service) FilterProductList(ctx context.Context, filterProductList TFilterProductList) ([]TFilteredProductList, error) {

	var filteredProductList []TFilteredProductList

	productList, err := s.productListRepo.FindByFilter(ctx, FilterProductListModel(filterProductList), []string{"name", "price", "discount", "imgLink", "_id", "detail", "batchId", "slug", "code"}, true)
	if err != nil {
		return nil, err
	}

	for _, item := range productList {
		filteredProductList = append(filteredProductList, TFilteredProductList{
			Name:     item.Name,
			Price:    item.Price,
			ImgLink:  item.ImgLink,
			Discount: item.Discount,
			Detail:   item.Detail.Hex(),
			Slug:     item.Slug,
			BatchId:  item.BatchId,
			Code:     item.Code,
			Id:       item.ID.Hex(),
		})
	}

	return filteredProductList, nil
}

func (s *Service) GetProductFilter(ctx context.Context, filterProductList TFilterProductList) (any, error) {
	res, err := s.productListRepo.GetFitler(ctx, FilterProductListModel(filterProductList), []string{"name", "price", "discount", "imgLink", "_id", "detail", "batchId", "slug", "code"}, true)

	if err != nil {
		return nil, err
	}

	finalResult := map[string]map[string]int{"color": map[string]int{}, "size": map[string]int{}}

	for _, item := range res {
		finalResult["color"][item.Color] += item.TotalStock

		for _, size := range item.Sizes {
			finalResult["size"][size] += item.TotalStock
		}
	}

	return finalResult, nil

}

func (s *Service) AddProductsToCollection(ctx context.Context, productData TAddProductToCollection) error {
	return s.productListRepo.AddProductsToCollection(ctx, AddProductToCollectionModel(productData))
}

func (s *Service) CreateProduct(ctx context.Context, productDetails TCreateProduct, creatorId string) error {

	_, err := s.txn.RunInTxn(ctx, func(sessCtx context.Context) (interface{}, error) {

		previewImg := productDetails.ProductDetail.ImgLink[0]
		productDetails.ProductList.ImgLink = previewImg

		productCode := strconv.Itoa(utils.GenerateRandomNumber(7))

		productDetailModal := CreateProductDetailModel{
			Code:        productCode,
			Name:        productDetails.Name,
			Slug:        productDetails.Slug,
			Description: productDetails.ProductDetail.Description,
			Variations:  productDetails.ProductDetail.Variations,
			ImgLink:     productDetails.ProductDetail.ImgLink,
			BatchId:     productDetails.BatchId,
		}

		productDetailModal.Variations = append(productDetailModal.Variations, Variation{Size: constant.BASE_SIZE, Price: productDetails.ProductList.Price, Discount: productDetails.ProductList.Discount, Stock: productDetails.ProductList.Stock})

		productDetailId, err := s.productDetailRepo.Create(ctx, productDetailModal, creatorId)
		if err != nil {
			return nil, err
		}
		sizes := []string{}

		for _, variation := range productDetails.ProductDetail.Variations {
			sizes = append(sizes, variation.Size)
		}

		productListModal := CreateProductListModel{
			Name:       productDetails.Name,
			Slug:       productDetails.Slug,
			Code:       productCode,
			ImgLink:    productDetails.ProductList.ImgLink,
			BatchId:    productDetails.BatchId,
			Sizes:      sizes,
			Color:      productDetails.ProductList.Color,
			Price:      productDetails.ProductList.Price,
			Stock:      productDetails.ProductList.Stock,
			Discount:   productDetails.ProductList.Discount,
			Detail:     productDetailId,
			Category:   productDetails.ProductList.Category,
			Gender:     productDetails.ProductList.Gender,
			Collection: productDetails.ProductList.Collection,
			Tags:       productDetails.ProductList.Tags,
		}
		productListId, err := s.productListRepo.Create(ctx, productListModal, creatorId)
		if err != nil {
			return nil, err
		}
		batchUpdatePayload := TBatchProductDetails{
			ImgLink:     productDetails.ProductList.ImgLink,
			ProductCode: productCode,
			Slug:        productDetails.Slug,
		}

		if err := s.productBatchRepo.UpdateBatchImg(ctx, productDetails.BatchId, batchUpdatePayload); err != nil {
			return "", err
		}

		return productListId, nil
	})

	return err
}

func (s *Service) GetVariations(ctx context.Context, productCode string) ([]Variation, error) {

	project := []string{"variations"}

	res, err := s.productDetailRepo.FindByCode(ctx, productCode, project, true)
	return res.Variations, err
}

func (s *Service) GetProductDetails(ctx context.Context, id string) (ProductDetail, error) {

	var result ProductDetail

	productDetails, err := s.productDetailRepo.FindByCode(ctx, id, []string{}, false)
	if err != nil {
		return result, err
	}
	return productDetails, nil
}

func (s *Service) GetProductBatchDetails(ctx context.Context, code string) (ProductBatch, error) {
	var result ProductBatch
	batchDetails, err := s.productBatchRepo.FindByCode(ctx, code, []string{}, false)
	if err != nil {
		return result, err
	}
	return batchDetails, nil
}

func (s *Service) CreateProductBatch(ctx context.Context, createPayload TCreateProductBatch, userId string) (string, error) {

	batchId := strconv.Itoa(utils.GenerateRandomNumber(6))

	payload := CreateProductBatchModel{
		Code:      batchId,
		Name:      createPayload.Name,
		CreatorId: userId,
	}

	return s.productBatchRepo.Create(ctx, payload)
}

func (s *Service) CreateProductStatus(ctx context.Context, createPayload CreateProductStatusModal, userId string) (string, error) {
	return s.productStatusRepo.Create(ctx, createPayload)
}

func (s *Service) GetProductRating(ctx context.Context, productCode string) (ProductStatus, error) {
	project := []string{"status", "rating"}
	return s.productStatusRepo.FindById(ctx, productCode, project, true)
}

// func (s *Service) UpdateProductRating(ctx context.Context, productCode string, rating int) {
// 	s.productStatusRepo.UpdateByProductCode(ctx, productCode, "rating", rating)
// }

func (s *Service) CreateProductTemplate(ctx context.Context, payload ProductTemplateModel) error {
	return s.productTemplateRepo.Create(ctx, payload)
}

func (s *Service) FindProductTemplate(ctx context.Context, payload FilterProductListModel) ([]ProductList, error) {

	project := []string{"name", "_id"}

	return s.productTemplateRepo.FindByFilter(ctx, payload, project, true)
}

func (s *Service) FindProductTemplateById(ctx context.Context, templateId string) (ProductTemplate, error) {
	return s.productTemplateRepo.FindById(ctx, templateId, []string{}, true)
}

func (s *Service) FindProductBatch(ctx context.Context, filterPayload FilterProductBatchListModel) ([]ProductBatch, error) {

	project := []string{"code", "name", "_id"}

	return s.productBatchRepo.FindByFilter(ctx, filterPayload, project, true)
}
