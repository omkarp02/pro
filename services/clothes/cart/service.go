package cart

import (
	"context"

	"github.com/omkarp02/pro/services/clothes/product"
	"github.com/omkarp02/pro/utils"
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

func (s *Service) AddToCard(ctx context.Context, userId string, cartDetails TAddToCart) error {

	var productCodes []string
	var sizes []string

	for _, item := range cartDetails.Items {
		productCodes = append(productCodes, item.ProductCode)
		sizes = append(sizes, item.Size)
	}

	sizes = utils.RemoveDuplicateStringFromSlice(sizes)
	productCodes = utils.RemoveDuplicateStringFromSlice(productCodes)
	project := []string{"code"}

	productDetailList, err := s.productRepo.GetProductsPriceBySizes(ctx, productCodes, sizes, project, true)

	if err != nil {
		return err
	}

	curCartTotalPrice := 0.0
	curCartTotalItems := 0

	for _, item := range cartDetails.Items {
		for _, product := range productDetailList {
			if item.ProductCode == product.Code {
				for _, variation := range product.Variations {
					if item.Size == variation.Size {
						curCartTotalPrice += float64(item.Quantity) * variation.Price
					}
				}
			}

		}
		curCartTotalItems += item.Quantity
	}

	newCartData := CreateCartModel{
		UserId:            userId,
		Items:             cartDetails.Items,
		CurCartTotalItems: curCartTotalItems,
		CurTotalPrice:     curCartTotalPrice,
	}

	err = s.repo.Create(ctx, newCartData)
	return err

}

func (s *Service) UpdateQuantityOfItem(ctx context.Context, userId string, payload IUpdateQuantityOfItem) error {

	return s.repo.UpdateCartItemQuantity(ctx, userId, payload.CartId, payload.Quantity)
}

func (s *Service) FindOne(ctx context.Context, userId string) (IFindOneRes, error) {

	var res IFindOneRes
	var productCodes []string
	var sizes []string

	cartDetails, err := s.repo.FindById(ctx, userId, []string{}, false)
	if err != nil {
		return res, err
	}

	for _, item := range cartDetails.Items {
		productCodes = append(productCodes, item.ProductCode)
		sizes = append(sizes, item.Size)
	}

	sizes = utils.RemoveDuplicateStringFromSlice(sizes)
	productCodes = utils.RemoveDuplicateStringFromSlice(productCodes)
	project := []string{"_id", "name", "imgLink", "code"}

	productDetailList, err := s.productRepo.GetProductsPriceBySizes(ctx, productCodes, sizes, project, true)

	if err != nil {
		return res, err
	}

	var items []IFindOneResCartItem

	for _, cartItem := range cartDetails.Items {
		for _, product := range productDetailList {
			if product.Code == cartItem.ProductCode {
				items = append(items, IFindOneResCartItem{
					CartId:      cartItem.CartId,
					ProductCode: cartItem.ProductCode,
					Size:        cartItem.Size,
					Quantity:    cartItem.Quantity,
					Product: IFindOneResProductItems{
						ID:         product.ID.Hex(),
						Name:       product.Name,
						PreviewImg: product.ImgLink[0],
						Variations: product.Variations,
					},
				})
			}
		}
	}

	res = IFindOneRes{
		ID:         cartDetails.ID.Hex(),
		Items:      items,
		TotalItems: cartDetails.TotalItems,
		TotalPrice: cartDetails.TotalPrice,
		Timestamps: cartDetails.Timestamps,
	}

	return res, nil

}
