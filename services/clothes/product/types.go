package product

type TCreateProductList struct {
	Name       string   `json:"name,omitempty" validate:"required"`
	Sizes      []string `json:"sizes,omitempty"`
	Color      string   `json:"color,omitempty" validate:"required"`
	Price      float64  `json:"price,omitempty" validate:"required"`
	ImgLink    string   `json:"imgLink,omitempty" validate:"required"`
	Stock      int      `json:"stock,omitempty" validate:"required"`
	Discount   int      `json:"discount,omitempty" validate:"required"`
	Detail     string   `json:"detail,omitempty"`
	Category   string   `json:"category,omitempty" validate:"required"`
	BatchId    string   `json:"batchId,omitempty" validate:"required"`
	Gender     string   `json:"gender,omitempty" validate:"required"`
	Collection []string `json:"collection,omitempty" validate:"required"`
	Tags       []string `json:"tags,omitempty" validate:"required"`
}

type TCreateProduct struct {
	ProductList   TCreateProductList   `json:"productList,omitempty" validate:"required"`
	ProductDetail TCreateProductDetail `json:"detail,omitempty" validate:"required"`
}

type TGetProductDetails struct {
	ProductId string `json:"productId,omitempty"`
}

type TCreateProductDetail struct {
	Name        string      `json:"name,omitempty"`
	PreviewImg  string      `json:"previewImg,omitempty"`
	Description Description `json:"description,omitempty"  validate:"required"`
	Variations  []Variation `json:"variations,omitempty"  validate:"required"`
	ImgLink     []string    `json:"imgLink,omitempty"  validate:"required"`
	BatchId     string      `json:"batchId,omitempty"  validate:"required"`
}

type TFilterProductList struct {
	Sizes      []string `query:"sizes,omitempty"`
	Color      string   `query:"color,omitempty"`
	MinPrice   float64  `query:"min_price,omitempty"`
	MaxPrice   float64  `query:"max_price,omitempty"`
	Collection string   `query:"collection,omitempty"`
	Name       string   `query:"name,omitempty"`
	Page       int      `query:"page,omitempty" validate:"required"`
	Limit      int      `query:"limit,omitempty" validate:"required"`
}

type TFilteredProductList struct {
	Id       string  `json:"id,omitempty"`
	Name     string  `json:"name,omitempty"`
	Price    float64 `json:"price,omitempty"`
	ImgLink  string  `json:"imgLink,omitempty"`
	Discount int     `json:"discount,omitempty"`
}

type TAddProductToCollection struct {
	CollectionName string   `json:"collection_name,omitempty" validate:"required"`
	ProductId      []string `json:"product_id,omitempty" validate:"required"`
}

type TBatchProductDetails struct {
	Id      string `json:"id,omitempty" bson:"id,omitempty"`
	ImgLink string `json:"imgLink,omitempty" bson:"imgLink,omitempty"`
}

type TProductDetailsServiceResponse struct {
	ProductDetails ProductDetail `json:"product_details,omitempty"`
	BatchDetails   ProductBatch  `json:"product_batch,omitempty"`
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
	Name        string      `json:"name,omitempty"`
	PreviewImg  string      `json:"previewImg,omitempty"`
	Description Description `json:"description,omitempty"`
	Variations  []Variation `json:"variations,omitempty"`
	ImgLink     []string    `json:"imgLink,omitempty"`
	BatchId     string      `json:"batchId,omitempty"`
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
	Sizes      []string `json:"sizes,omitempty"`
	Color      string   `json:"color,omitempty"`
	MinPrice   float64  `json:"min_price,omitempty"`
	MaxPrice   float64  `json:"max_price,omitempty"`
	Collection string   `json:"collection,omitempty"`
	Name       string   `json:"name,omitempty"`
	Page       int      `json:"page,omitempty"`
	Limit      int      `json:"limit,omitempty"`
}

type CreateProductBatchModel struct {
	BatchCode   string                `json:"batchCode,omitempty"`
	ProductList []BatchProductDetails `json:"batchProductDetails,omitempty"`
}

type AddProductToCollectionModel struct {
	CollectionName string   `json:"collection_name,omitempty"`
	ProductId      []string `json:"product_id,omitempty"`
}
