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

	productList, err := s.productListRepo.FindByFilter(ctx, FilterProductListModel(filterProductList), []string{"name", "price", "discount", "imgLink", "_id", "detail", "batchId"}, true)
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
			BatchId:  item.BatchId,
			Id:       item.ID.Hex(),
		})
	}

	return filteredProductList, nil
}

func (s *Service) AddProductsToCollection(ctx context.Context, productData TAddProductToCollection) error {
	return s.productListRepo.AddProductsToCollection(ctx, AddProductToCollectionModel(productData))
}

func (s *Service) CreateProduct(ctx context.Context, productDetails TCreateProduct) error {

	_, err := s.txn.RunInTxn(ctx, func(sessCtx context.Context) (interface{}, error) {

		previewImg := productDetails.ProductDetail.ImgLink[0]
		productDetails.ProductList.ImgLink = previewImg
		productDetails.ProductDetail.PreviewImg = previewImg
		productDetails.ProductDetail.Name = productDetails.ProductList.Name
		productDetails.ProductDetail.Variations = append(productDetails.ProductDetail.Variations, Variation{Size: constant.BASE_SIZE, Price: productDetails.ProductList.Price, Discount: productDetails.ProductList.Discount})

		productDetailId, err := s.productDetailRepo.Create(ctx, CreateProductDetailModel(productDetails.ProductDetail))
		if err != nil {
			return nil, err
		}

		sizes := []string{}

		for _, variation := range productDetails.ProductDetail.Variations {
			sizes = append(sizes, variation.Size)
		}

		productDetails.ProductList.Detail = productDetailId
		productDetails.ProductList.Sizes = sizes

		productListId, err := s.productListRepo.Create(ctx, CreateProductListModel(productDetails.ProductList))
		if err != nil {
			return nil, err
		}

		batchUpdatePayload := TBatchProductDetails{
			ImgLink:         productDetails.ProductList.ImgLink,
			ProductListId:   productListId,
			ProductDetailId: productDetailId,
		}

		if err := s.productBatchRepo.UpdateBatchImg(ctx, productDetails.ProductList.BatchId, batchUpdatePayload); err != nil {
			return "", err
		}

		return productListId, nil
	})

	return err
}

func (s *Service) GetProductDetails(ctx context.Context, productId string) (ProductDetail, error) {

	var result ProductDetail

	productDetails, err := s.productDetailRepo.FindById(ctx, productId, []string{}, false)
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
