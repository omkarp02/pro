package categories

import "github.com/omkarp02/pro/services/utils/store"

type TCreateCategory struct {
	Name        string `json:"name,omitempty" validate:"required"`
	Description string `json:"description,omitempty" validate:"required"`
	ImgLink     string `json:"imgLink,omitempty" validate:"required"`
	Icon        string `json:"icon,omitempty" validate:"required"`
	Status      string `json:"status,omitempty" validate:"required,validStatus"`
	Slug        string `json:"slug,omitempty" validate:"required"`
}

type CreateCategoryModal struct {
	CatId string `json:"catId,omitempty"`
	TCreateCategory
}

type TFilterCategory struct {
	Pagination store.Pagination `json:"pagination,inline" validation:"required"`
}

type TCategoryList struct {
	Id    string `json:"id,omitempty"`
	CatId string `json:"catId,omitempty"`
	Name  string `json:"name,omitempty"`
	Icon  string `json:"icon,omitempty"`
	Slug  string `json:"slug,omitempty"`
}

type FilterCategoryModal struct {
	Status     string
	Pagination store.Pagination `json:"pagination,inline"`
}
