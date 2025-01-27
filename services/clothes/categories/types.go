package categories

import "github.com/omkarp02/pro/services/utils/store"

type TCreateCategory struct {
	Name        string `json:"name,omitempty" validate:"required"`
	Description string `json:"description,omitempty" validate:"required"`
	ImgLink     string `json:"imgLink,omitempty" validate:"required"`
	Icon        string `json:"icon,omitempty" validate:"required"`
	IsActive    bool   `json:"isActive,omitempty" validate:"required"`
	Slug        string `json:"slug,omitempty" validate:"required"`
}

type CreateCategoryModal struct {
	CatId string `json:"catId,omitempty"`
	TCreateCategory
}

type TFilterCategory struct {
	IsActive   bool             `json:"isActive,omitempty"`
	Pagination store.Pagination `json:"pagination,inline" validation:"required"`
}

type TCategoryList struct {
	CatId string `json:"catId,omitempty"`
	Name  string `json:"name,omitempty"`
	Icon  string `json:"icon,omitempty"`
	Slug  string `json:"slug,omitempty"`
}

type FilterCategoryModal struct {
	IsActive   bool             `json:"isActive,omitempty"`
	Pagination store.Pagination `json:"pagination,inline"`
}
