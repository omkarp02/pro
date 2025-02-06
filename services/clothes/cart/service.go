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

	// var productCodes []string
	// var sizes []string

	// for _, item := range cartDetails.Items {
	// 	productCodes = append(productCodes, item.ProductCode)
	// 	sizes = append(sizes, item.Size)
	// }

	// sizes = utils.RemoveDuplicateStringFromSlice(sizes)
	// productCodes = utils.RemoveDuplicateStringFromSlice(productCodes)
	// project := []string{"code"}

	// productDetailList, err := s.productRepo.GetProductsPriceBySizes(ctx, productCodes, sizes, project, true)

	// if err != nil {
	// 	return err
	// }

	// curCartTotalPrice := 0.0
	// curCartTotalItems := 0

	// for _, item := range cartDetails.Items {
	// 	for _, product := range productDetailList {
	// 		if item.ProductCode == product.Code {
	// 			for _, variation := range product.Variations {
	// 				if item.Size == variation.Size {
	// 					curCartTotalPrice += float64(item.Quantity) * variation.Price
	// 				}
	// 			}
	// 		}

	// 	}
	// 	curCartTotalItems += item.Quantity
	// }

	newCartData := CreateCartModel{
		UserId: userId,
		Item:   cartDetails.Item,
		// CurCartTotalItems: curCartTotalItems,
		// CurTotalPrice:     curCartTotalPrice,
	}

	err := s.repo.Create(ctx, newCartData)
	return err

}

func (s *Service) UpdateQuantityOfItem(ctx context.Context, userId string, payload IUpdateQuantityOfItem) error {
	field := "quantity"
	var value interface{} = payload.Quantity
	if len(payload.Size) > 0 {
		field = "size"
		value = payload.Size
	}

	return s.repo.UpdateCartItem(ctx, userId, payload.CartId, field, value)
}

func (s *Service) FindOne(ctx context.Context, userId string) (IFindOneRes, error) {

	var res IFindOneRes

	cartDetails, err := s.repo.FindById(ctx, userId, []string{}, false)
	if err != nil {
		return res, err
	}

	items, err := s.PopulateProductDetailsInCartItem(ctx, cartDetails.Items)
	if err != nil {
		return res, err
	}

	res = IFindOneRes{
		ID:    cartDetails.ID.Hex(),
		Items: items,
		// TotalItems: cartDetails.TotalItems,
		// TotalPrice: cartDetails.TotalPrice,
		Timestamps: cartDetails.Timestamps,
	}

	return res, nil

}

func (s *Service) PopulateProductDetailsInCartItem(ctx context.Context, cartItems []CartItem) ([]IFindOneResCartItem, error) {

	var items []IFindOneResCartItem

	if len(cartItems) == 0 {
		return items, nil
	}

	var productCodes []string
	var sizes []string
	for _, item := range cartItems {
		productCodes = append(productCodes, item.ProductCode)
		sizes = append(sizes, item.Size)
	}

	sizes = utils.RemoveDuplicateStringFromSlice(sizes)
	productCodes = utils.RemoveDuplicateStringFromSlice(productCodes)
	project := []string{"_id", "name", "imgLink", "code"}

	productDetailList, err := s.productRepo.GetProductsPriceBySizes(ctx, productCodes, sizes, project, true)

	if err != nil {
		return items, err
	}

	for _, cartItem := range cartItems {
		for _, _product := range productDetailList {
			if _product.Code == cartItem.ProductCode {
				var variation product.Variation
				for _, varItem := range _product.Variations {
					if varItem.Size == cartItem.Size {
						variation = varItem
						break
					}
				}
				items = append(items, IFindOneResCartItem{
					CartId:      cartItem.CartId,
					ProductCode: cartItem.ProductCode,
					Size:        cartItem.Size,
					Quantity:    cartItem.Quantity,
					Product: IFindOneResProductItems{
						ID:         _product.ID.Hex(),
						Name:       _product.Name,
						PreviewImg: _product.ImgLink[0],
						Variations: variation,
					},
				})
				break
			}
		}
	}

	return items, nil
}

func (s *Service) GetCartItemForOffline(ctx context.Context, cartDetails GetCartOfflineModal) ([]IFindOneResCartItem, error) {
	length := len(cartDetails.ProductCode)
	cartItems := make([]CartItem, length)

	for i := 0; i < length; i++ {
		cartItems[i] = CartItem{
			ProductCode: cartDetails.ProductCode[i],
			Size:        cartDetails.Size[i],
			CartId:      cartDetails.CartId[i],
		}
	}

	res, err := s.PopulateProductDetailsInCartItem(ctx, cartItems)
	return res, err
}

func (s *Service) DeleteCartItem(ctx context.Context, cartId string, userId string) error {

	return s.repo.DeleteCartItem(ctx, userId, cartId)
}

func (s *Service) GetTotalItems(ctx context.Context, userId string) (int, error) {

	project := []string{}

	res, err := s.repo.FindById(ctx, userId, project, false)
	if err != nil {
		return 0, err
	}
	totalItem := 0

	for _, item := range res.Items {
		totalItem += item.Quantity
	}
	return totalItem, nil

}
