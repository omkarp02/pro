package product

type TVariation struct {
	Size  string  `json:"size,omitempty"`
	Price float64 `json:"price,omitempty"`
}

type TPrice struct {
	BasePrice  float64      `json:"basePrice,omitempty"`
	Variations []TVariation `json:"variations,omitempty"`
}

type TCreateProductList struct {
	Name       string   `json:"name,omitempty" validate:"required"`
	Sizes      []string `json:"sizes,omitempty" validate:"required,min=1,max=10"`
	Color      string   `json:"color,omitempty" validate:"required"`
	Price      float64  `json:"price,omitempty" validate:"required"`
	ImgLink    string   `json:"imgLink,omitempty" validate:"required"`
	Stock      int      `json:"stock,omitempty" validate:"required"`
	Discount   int      `json:"discount,omitempty" validate:"required"`
	Detail     string   `json:"detail,omitempty" validate:"required"`
	Category   string   `json:"category,omitempty" validate:"required"`
	BatchId    string   `json:"batchId,omitempty" validate:"required"`
	Gender     string   `json:"gender,omitempty" validate:"required"`
	Collection []string `json:"collection,omitempty" validate:"required"`
	Tags       []string `json:"tags,omitempty" validate:"required"`
}

type TCreateProductDetail struct {
	Description Description `json:"description,omitempty"`
	Variations  []Variation `json:"variations,omitempty"`
}

type TFilterProductList struct {
	Sizes    []string `query:"sizes,omitempty"`
	Color    string   `query:"color,omitempty"`
	MinPrice float64  `query:"min_price,omitempty"`
	MaxPrice float64  `query:"max_price,omitempty"`
	Name     string   `query:"name,omitempty"`
	Page     int      `query:"page,omitempty" validate:"required"`
	Limit    int      `query:"limit,omitempty" validate:"required"`
}

type TAddProductToCollection struct {
	CollectionName string   `json:"collection_name,omitempty" validate:"required"`
	ProductId      []string `json:"product_id,omitempty" validate:"required"`
}

// here are the model types
type VariationModel struct {
	Size  string  `json:"size,omitempty"`
	Price float64 `json:"price,omitempty"`
}

type PriceModel struct {
	BasePrice  float64     `json:"basePrice,omitempty"`
	Variations []Variation `json:"variations,omitempty"`
}

type CreateProductDetailModel struct {
	Description Description `json:"description,omitempty"`
	Variations  []Variation `json:"variations,omitempty"`
}
type CreateProductListModel struct {
	Name       string   `json:"name,omitempty"`
	Sizes      []string `json:"sizes,omitempty"`
	Color      string   `json:"color,omitempty"`
	Price      float64  `json:"price,omitempty"`
	ImgLink    string   `json:"imgLink,omitempty"`
	Stock      int      `json:"stock,omitempty"`
	Discount   int      `json:"discount,omitempty"`
	Detail     string   `json:"detail,omitempty"`
	Category   string   `json:"category,omitempty"`
	BatchId    string   `json:"batchId,omitempty"`
	Gender     string   `json:"gender,omitempty"`
	Collection []string `json:"collection,omitempty"`
	Tags       []string `json:"tags,omitempty"`
}

type FilterProductListModel struct {
	Sizes    []string `json:"sizes,omitempty"`
	Color    string   `json:"color,omitempty"`
	MinPrice float64  `json:"min_price,omitempty"`
	MaxPrice float64  `json:"max_price,omitempty"`
	Name     string   `json:"name,omitempty"`
	Page     int      `json:"page,omitempty"`
	Limit    int      `json:"limit,omitempty"`
}

type AddProductToCollectionModel struct {
	CollectionName string   `json:"collection_name,omitempty"`
	ProductId      []string `json:"product_id,omitempty"`
}
