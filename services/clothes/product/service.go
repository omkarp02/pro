package product

import "context"

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

func (s *Service) CreateProductList(ctx context.Context, createProductList TCreateProductList) (string, error) {
	return s.productListRepo.Create(ctx, CreateProductListModel(createProductList))
}

func (s *Service) FilterProductList(ctx context.Context, filterProductList TFilterProductList) ([]ProductList, error) {
	return s.productListRepo.FindByFilter(ctx, FilterProductListModel(filterProductList))
}

func (s *Service) AddProductsToCollection(ctx context.Context, productData TAddProductToCollection) error {
	return s.productListRepo.AddProductsToCollection(ctx, AddProductToCollectionModel(productData))
}
