package filter

type TCreateFilter struct {
	Name     string `json:"name,omitempty" validate:"required"`
	Type     string `json:"type,omitempty" validate:"required"`
	Category string `json:"category,omitempty" validate:"required"`
	Slug     string `json:"slug,omitempty"  validate:"required"`
}

type TCreateFilterType struct {
	Name string `json:"name,omitempty" validate:"required"`
}

type TFilterTypeList struct {
	Page  int `query:"page,omitempty" validate:"required"`
	Limit int `query:"limit,omitempty" validate:"required"`
}

type TFilterList struct {
	Page     int    `query:"page,omitempty"`
	Limit    int    `query:"limit,omitempty"`
	Category string `query:"category,omitempty"`
	Type     string `query:"type,omitempty"`
}

type CreateFilterModal struct {
	Name     string `json:"name,omitempty"`
	Type     string `json:"type,omitempty"`
	Category string `json:"category,omitempty"`
	Slug     string `json:"slug,omitempty"`
	Status   string `json:"status,omitempty"`
}
type CreateFilterTypeModal struct {
	Name   string `json:"name,omitempty"`
	Status string `json:"status,omitempty"`
}

type FilterItemForUserRes struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}

type FilterListForUserRes struct {
	Name   string                 `json:"name,omitempty"`
	Fitler []FilterItemForUserRes `json:"filter,omitempty"`
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
