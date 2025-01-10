package order

import "github.com/omkarp02/pro/services/utils/store"

type TCreateOrder struct {
	Address store.TAddress `json:"address,omitempty"`
}

type CreateOrderItemModel struct {
	ProductId string  `json:"productId,omitempty"`
	Price     float64 `json:"price,omitempty"`
	Size      string  `json:"size,omitempty"`
	Quantity  int     `json:"quantity,omitempty"`
	ItemId    string  `json:"itemId,omitempty"`
	OrderId   string  `json:"orderId,omitempty"`
	UserId    string  `json:"userId,omitempty"`
	Status    string  `json:"status,omitempty"`
}

type TShippingInfo struct {
	Address store.AddressModel `json:"address,omitempty"`
}

type CreateOrderModal struct {
	UserId        string        `json:"userId,omitempty"`
	ShippingInfo  TShippingInfo `json:"shippingInfo,omitempty"`
	Items         []string      `json:"items,omitempty"`
	TotalPrice    float64       `json:"totalPrice,omitempty"`
	ShippingPrice float64       `json:"shippingPrice,omitempty"`
}
