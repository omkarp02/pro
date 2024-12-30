package categories

type TCreateCategory struct {
	Name string `json:"name,omitempty" validate:"required"`
}

type CreateCategoryModal struct {
	Name string `json:"name,omitempty" bson:"name"`
}
