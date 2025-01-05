package cart

type TCartItem struct {
	ProductId string `json:"productId,omitempty"  validate:"required"` // ID of the product
	Size      string `json:"size,omitempty"  validate:"required"`
	Quantity  int    `json:"quantity,omitempty"  validate:"required"` // Quantity of the product
	CartId    string `json:"cartId,omitempty"`
}

type TAddToCart struct {
	UserId string    `json:"userId,omitempty"  validate:"required"`
	Items  TCartItem `json:"item,omitempty"  validate:"required"`
}

type CreateCartModel struct {
	UserId string    `json:"userId,omitempty"` // User owning the cart
	Items  TCartItem `json:"items,omitempty"`  // List of items in the cart
	Price  float64
}
