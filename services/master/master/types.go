package master

type TCreate struct {
	Type string `json:"type,omitempty" validate:"required,oneof=bussiness-category gender color product-collection"`
	Name string `json:"name,omitempty" validate:"required"`
}

type CreateModal struct {
	TCreate
	CreatorId string
}

type TFilterList struct {
	Page  int      `query:"page,omitempty" validate:"required"`
	Limit int      `query:"limit,omitempty" validate:"required"`
	Types []string `query:"types,omitempty" validate:"required,dive,oneof=bussiness-category gender color product-collection"`
}

type FilterListModel struct {
	Page  int
	Limit int
	Types []string
}
