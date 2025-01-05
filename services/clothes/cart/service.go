package cart

import (
	"context"

	"github.com/omkarp02/pro/services/clothes/product"
)

type Service struct {
	repo        *Repo
	productRepo *product.ProductDetailRepo
}

func NewService(repo *Repo, productRepo *product.ProductDetailRepo) *Service {
	return &Service{
		repo:        repo,
		productRepo: productRepo,
	}
}

func (s *Service) AddToCard(ctx context.Context, cartDetails TAddToCart) error {

	productDetailPrice, err := s.productRepo.GetProductPriceBySize(ctx, cartDetails.Items.ProductId, cartDetails.Items.Size)

	if err != nil {
		return err
	}

	newCartData := CreateCartModel{
		UserId: cartDetails.UserId,
		Items:  cartDetails.Items,
		Price:  productDetailPrice,
	}

	err = s.repo.Create(ctx, newCartData)
	return err

}
