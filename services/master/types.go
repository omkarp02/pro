package master

type TCreate struct {
	Type string `json:"type,omitempty" validate:"required,oneof=state country"`
	Name string `json:"name,omitempty" validate:"required"`
}

type CreateModal struct {
	TCreate
	CreatorId string
}

type TFilterList struct {
	Page  int    `query:"page,omitempty" validate:"required"`
	Limit int    `query:"limit,omitempty" validate:"required"`
	Type  string `json:"type,omitempty" validate:"required,oneof=state country"`
}

type FilterListModel struct {
	Page  int
	Limit int
	Type  string
}
