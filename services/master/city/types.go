package city

type TCreate struct {
	Name    string `json:"name,omitempty" validate:"required"`
	StateId string `json:"stateId,omitempty" validate:"required"`
}

type TFilterList struct {
	StateId string `query:"stateId,omitempty"`
	Name    string `query:"name,omitempty"`
	Page    int    `query:"page,omitempty" validate:"required"`
	Limit   int    `query:"limit,omitempty" validate:"required"`
}

// here this are type of repo which are entire seperately handled
type CreateModal struct {
	Name    string
	StateId string
}

type FilterListModel struct {
	StateId string `json:"stateId,omitempty"`
	Name    string `json:"name,omitempty"`
	Page    int    `json:"page,omitempty"`
	Limit   int    `json:"limit,omitempty"`
}
