package order

import (
	"context"
	"strconv"

	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/services/clothes/cart"
	"github.com/omkarp02/pro/services/clothes/product"
	"github.com/omkarp02/pro/services/utils/store"
	"github.com/omkarp02/pro/utils"
	"github.com/omkarp02/pro/utils/errutil"
)

type Service struct {
	repo          *Repo
	orderItemRepo *OrderItemRepo
	cartRepo      *cart.Repo
	productRepo   *product.ProductDetailRepo
	txn           db.TransactionManager
}

func NewService(repo *Repo, orderItemRepo *OrderItemRepo, cartRepo *cart.Repo, productRepo *product.ProductDetailRepo, txn db.TransactionManager) *Service {
	return &Service{
		repo:          repo,
		orderItemRepo: orderItemRepo,
		cartRepo:      cartRepo,
		productRepo:   productRepo,
		txn:           txn,
	}
}

func (s *Service) CreateOrder(ctx context.Context, userId string, createOrderPayload TCreateOrder) (string, error) {
	var productCodes []string
	var sizes []string

	cartDetails, err := s.cartRepo.FindById(ctx, userId, []string{}, false)
	if err != nil {
		return "", err
	}

	for _, item := range cartDetails.Items {
		productCodes = append(productCodes, item.ProductCode)
		sizes = append(sizes, item.Size)
	}

	sizes = utils.RemoveDuplicateStringFromSlice(sizes)
	productCodes = utils.RemoveDuplicateStringFromSlice(productCodes)
	project := []string{"_id", "name", "previewImg", "code"}

	productDetailList, err := s.productRepo.GetProductsPriceBySizes(ctx, productCodes, sizes, project, true)

	if err != nil {
		return "", err
	}

	if len(cartDetails.Items) == 0 {
		return "", errutil.InternalServerError("Cart is empty")
	}

	orderItemList, totalPrice, err := createOrderItem(userId, cartDetails, productDetailList)
	if err != nil {
		return "", err
	}

	result, err := s.txn.RunInTxn(ctx, func(sessCtx context.Context) (interface{}, error) {
		ids, err := s.orderItemRepo.InsertMany(ctx, orderItemList)
		if err != nil {
			return "", err
		}

		orderItem := CreateOrderModal{
			UserId: userId,
			ShippingInfo: TShippingInfo{
				Address: store.AddressModel(createOrderPayload.Address),
			},
			Items:         ids,
			TotalPrice:    totalPrice,
			ShippingPrice: 50,
		}

		return s.repo.Create(ctx, orderItem)
	})

	data, _ := result.(string)

	return data, err
}

func createOrderItem(userId string, cartDetails cart.Cart, productDetailList []product.ProductDetail) ([]CreateOrderItemModel, float64, error) {
	var orderItemList []CreateOrderItemModel

	productMap := make(map[string]product.ProductDetail, len(productDetailList))

	for _, product := range productDetailList {
		productMap[product.Code] = product
	}

	//here we are checking if product are in stock
	for _, cartItem := range cartDetails.Items {
		_, exists := productMap[cartItem.ProductCode]
		if !exists {
			return nil, 0.0, errutil.InternalServerError("Something went wrong")
		}
	}

	var totalProductPrice float64

	for _, cartItem := range cartDetails.Items {
		product := productMap[cartItem.ProductCode]
		var size string
		var price float64

		for _, variation := range product.Variations {
			if variation.Size == cartItem.Size {
				price = variation.Price
				break
			}
		}

		for i := 0; i < cartItem.Quantity; i++ {
			orderItem := CreateOrderItemModel{
				UserId:    userId,
				Price:     price,
				ProductId: product.ID.Hex(),
				Size:      size,
				Quantity:  cartItem.Quantity,
				ItemId:    "I-" + strconv.Itoa(utils.GenerateRandomNumber(5)),
				OrderId:   "dummyorderID",
				Status:    "dummySTatus",
			}

			totalProductPrice += price

			orderItemList = append(orderItemList, orderItem)
		}
	}

	return orderItemList, totalProductPrice, nil
}
