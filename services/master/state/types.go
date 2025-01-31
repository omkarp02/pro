package state

type TCreate struct {
	Name string `json:"name,omitempty" validate:"required"`
}

type TFilterList struct {
	Name  string `query:"name,omitempty"`
	Page  int    `query:"page,omitempty" validate:"required"`
	Limit int    `query:"limit,omitempty" validate:"required"`
}

// here this are type of repo which are entire seperately handled
type CreateModal struct {
	Name string `json:"name,omitempty"`
}

type FilterListModel struct {
	Name  string `json:"name,omitempty"`
	Page  int    `json:"page,omitempty"`
	Limit int    `json:"limit,omitempty"`
}
