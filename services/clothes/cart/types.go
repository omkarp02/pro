package cart

import (
	"github.com/omkarp02/pro/services/clothes/product"
	"github.com/omkarp02/pro/services/utils/store"
)

type TCartItem struct {
	ProductCode string `json:"productCode,omitempty"  validate:"required"`
	Size        string `json:"size,omitempty"  validate:"required"`
	Quantity    int    `json:"quantity,omitempty"  validate:"required"`
}

type IUpdateQuantityOfItem struct {
	Quantity int    `json:"quantity,omitempty"  validate:"required"`
	CartId   string `json:"cartId,omitempty"  validate:"required"`
}

type IFindOneResProductItems struct {
	ID         string              `json:"id,omitempty"`
	Name       string              `json:"name,omitempty" bson:"name,omitempty"`
	PreviewImg string              `json:"previewImg,omitempty" bson:"previewImg,omitempty"`
	Variations []product.Variation `json:"variations,omitempty" bson:"variations,omitempty"`
}

type IFindOneResCartItem struct {
	CartId      string                  `json:"cartId,omitempty"`
	ProductCode string                  `json:"productCode,omitempty"`
	Size        string                  `json:"size,omitempty"`
	Quantity    int                     `json:"quantity,omitempty"`
	Product     IFindOneResProductItems `json:"product,omitempty"`
}

type IFindOneRes struct {
	ID         string                `json:"_id,omitempty"`
	Items      []IFindOneResCartItem `json:"items,omitempty"`
	TotalItems int                   `bson:"totalItems,omitempty" json:"totalItems,omitempty"`
	TotalPrice float64               `bson:"totalPrice,omitempty" json:"totalPrice,omitempty"`
	Timestamps store.Timestamps      `bson:",inline" json:"timestamp"`
}

type TAddToCart struct {
	Items []TCartItem `json:"item,omitempty"  validate:"required,dive"`
}

type CreateCartModel struct {
	UserId            string      `json:"userId,omitempty"` // User owning the cart
	Items             []TCartItem `json:"items,omitempty"`  // List of items in the cart
	CurTotalPrice     float64
	CurCartTotalItems int
}
