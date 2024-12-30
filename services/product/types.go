package product

type CreateProductListModal struct {
	Name           string  `json:"name,omitempty" bson:"name,omitempty"`
	Size           string  `json:"size,omitempty" bson:"size,omitempty"`
	Color          string  `json:"color,omitempty" bson:"color,omitempty"`
	Brand          string  `json:"brand,omitempty" bson:"brand,omitempty"`
	Price          float64 `json:"price,omitempty" bson:"price,omitempty"`
	Discount       int     `json:"discount,omitempty" bson:"discount,omitempty"`
	Stock          int     `json:"stock,omitempty" bson:"stock,omitempty"`
	PreviewImgLink string  `json:"previewImgLink,omitempty" bson:"previewImgLink,omitempty"`
	ImgLink        string  `json:"imgLink,omitempty" bson:"imgLink,omitempty"`
}
