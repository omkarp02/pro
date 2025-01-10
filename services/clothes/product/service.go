package product

import (
	"context"
	"strconv"

	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/utils"
)

type Service struct {
	productListRepo   *ProductListRepo
	productDetailRepo *ProductDetailRepo
	productBatchRepo  *ProductBatchRepo
	txn               db.TransactionManager
}

func NewService(productListRepo *ProductListRepo, productDetailRepo *ProductDetailRepo, productBatchRepo *ProductBatchRepo, txn db.TransactionManager) *Service {
	return &Service{
		productListRepo:   productListRepo,
		productDetailRepo: productDetailRepo,
		productBatchRepo:  productBatchRepo,
		txn:               txn,
	}
}

func (s *Service) FilterProductList(ctx context.Context, filterProductList TFilterProductList) ([]TFilteredProductList, error) {

	var filteredProductList []TFilteredProductList

	productList, err := s.productListRepo.FindByFilter(ctx, FilterProductListModel(filterProductList), []string{"name", "price", "discount", "imgLink", "_id"}, true)
	if err != nil {
		return nil, err
	}

	for _, item := range productList {
		filteredProductList = append(filteredProductList, TFilteredProductList{
			Name:     item.Name,
			Price:    item.Price,
			ImgLink:  item.ImgLink,
			Discount: item.Discount,
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

		productDetails.ProductDetail.PreviewImg = productDetails.ProductList.ImgLink
		productDetails.ProductDetail.Name = productDetails.ProductList.Name

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
			ImgLink: productDetails.ProductList.ImgLink,
			Id:      productListId,
		}

		if err := s.productBatchRepo.UpdateBatchImg(ctx, productDetails.ProductList.BatchId, batchUpdatePayload); err != nil {
			return "", err
		}

		return productListId, nil
	})

	return err
}

func (s *Service) GetProductDetails(ctx context.Context, productId string) (TProductDetailsServiceResponse, error) {

	var result TProductDetailsServiceResponse

	productDetails, err := s.productDetailRepo.FindById(ctx, productId, []string{}, false)
	if err != nil {
		return result, err
	}

	batchDetails, err := s.productBatchRepo.FindById(ctx, productDetails.BatchId.Hex(), []string{}, false)
	if err != nil {
		return result, err
	}

	//here need to link the batch details and also need to update the code on create like on product crate also update the batch document
	result.ProductDetails = productDetails
	result.BatchDetails = batchDetails

	return result, nil
}

func (s *Service) CreateProductBatch(ctx context.Context) (string, error) {

	batchId := "BATCH" + strconv.Itoa(utils.GenerateRandomNumber(5))

	payload := CreateProductBatchModel{
		BatchCode: batchId,
	}

	return s.productBatchRepo.Create(ctx, payload)
}
