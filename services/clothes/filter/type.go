package filter

type TCreateFilter struct {
	Name     string `json:"name,omitempty" validate:"required"`
	Type     string `json:"type,omitempty" validate:"required"`
	Category string `json:"category,omitempty" validate:"required"`
}

type TCreateFilterType struct {
	Name string `json:"name,omitempty" validate:"required"`
}

type TFilterTypeList struct {
	Page  int `query:"page,omitempty" validate:"required"`
	Limit int `query:"limit,omitempty" validate:"required"`
}

type TFilterList struct {
	Page     int    `query:"page,omitempty" validate:"required"`
	Limit    int    `query:"limit,omitempty" validate:"required"`
	Category string `query:"category,omitempty"`
	Type     string `query:"type,omitempty"`
}

type CreateFilterModal struct {
	Name     string `json:"name,omitempty"`
	Type     string `json:"type,omitempty"`
	Category string `json:"category,omitempty"`
	Status   string `json:"status,omitempty"`
}
type CreateFilterTypeModal struct {
	Name   string `json:"name,omitempty"`
	Status string `json:"status,omitempty"`
}

type FilterTypeListModel struct {
	Page  int
	Limit int
}

type FilterListModel struct {
	Page     int
	Limit    int
	Category string
	Type     string
}
