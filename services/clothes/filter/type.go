package filter

type TCreateFilter struct {
	Name     string `json:"name,omitempty" validate:"required"`
	Type     string `json:"type,omitempty" validate:"required"`
	Category string `json:"category,omitempty" validate:"required"`
}

type TCreateFilterType struct {
	Name string `json:"name,omitempty"`
}

type CreateFilterModal struct {
	Name     string `json:"name,omitempty" bson:"name,omitempty"`
	Type     string `json:"type,omitempty" bson:"type,omitempty"`
	Category string `json:"category,omitempty" bson:"category,omitempty"`
}
type CreateFilterTypeModal struct {
	Name string `json:"name,omitempty"`
}
