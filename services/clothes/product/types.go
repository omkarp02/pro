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
	Name     string   `json:"name,omitempty"`
	Sizes    []string `json:"sizes,omitempty"`
	Color    string   `json:"color,omitempty"`
	Price    TPrice   `json:"price,omitempty"`
	ImgLink  string   `json:"imgLink,omitempty"`
	Stock    int      `json:"stock,omitempty"`
	Discount int      `json:"discount,omitempty"`
	Detail   string   `json:"detail,omitempty"`
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

type CreateProductListModel struct {
	Name     string   `json:"name,omitempty"`
	Sizes    []string `json:"sizes,omitempty"`
	Color    string   `json:"color,omitempty"`
	Price    Price    `json:"price,omitempty"`
	ImgLink  string   `json:"imgLink,omitempty"`
	Stock    int      `json:"stock,omitempty"`
	Discount int      `json:"discount,omitempty"`
	Detail   string   `json:"detail,omitempty"`
}
