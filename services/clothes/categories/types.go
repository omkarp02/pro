package categories

type TCreateCategory struct {
	CatId       string `json:"catId,omitempty" validate:"required"`
	Name        string `json:"name,omitempty" validate:"required"`
	Description string `json:"description,omitempty" validate:"required"`
	ImgLink     string `json:"imgLink,omitempty" validate:"required"`
	Icon        string `json:"icon,omitempty" validate:"required"`
	IsActive    bool   `json:"isActive,omitempty" validate:"required"`
	Slug        string `json:"slug,omitempty" validate:"required"`
}

type CreateCategoryModal struct {
	CatId       string `json:"catId,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	ImgLink     string `json:"imgLink,omitempty"`
	Icon        string `json:"icon,omitempty"`
	IsActive    bool   `json:"isActive,omitempty"`
	Slug        string `json:"slug,omitempty"`
}
