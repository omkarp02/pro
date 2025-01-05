package product

import (
	"context"
	"fmt"
)

type Service struct {
	productListRepo   *ProductListRepo
	productDetailRepo *ProductDetailRepo
}

func NewService(productListRepo *ProductListRepo, productDetailRepo *ProductDetailRepo) *Service {
	return &Service{
		productListRepo:   productListRepo,
		productDetailRepo: productDetailRepo,
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

	productDetailId, err := s.productDetailRepo.Create(ctx, CreateProductDetailModel(productDetails.ProductDetail))
	if err != nil {
		return err
	}
	sizes := []string{}

	for _, variation := range productDetails.ProductDetail.Variations {
		sizes = append(sizes, variation.Size)
	}

	productDetails.ProductList.Detail = productDetailId
	productDetails.ProductList.Sizes = sizes

	productListId, err := s.productListRepo.Create(ctx, CreateProductListModel(productDetails.ProductList))
	if err != nil {
		return err
	}

	fmt.Println(productListId)

	return nil
}
