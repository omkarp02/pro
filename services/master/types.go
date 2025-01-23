package master

type TCreate struct {
	Type  string `json:"type,omitempty" validate:"required,oneof=state country"`
	Label string `json:"label,omitempty" validate:"required"`
	Value string `json:"value,omitempty" validate:"required"`
}

type CreateModal struct {
	TCreate
	CreatorId string
}
