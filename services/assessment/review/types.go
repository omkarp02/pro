package review

type TCreate struct {
	ProductCode string `json:"productCode" validate:"required"`
	Rating      int    `json:"rating" validate:"required,min=1,max=5"`
	ReviewText  string `json:"reviewText" validate:"required"`
}

type TFilterList struct {
	Page  int `query:"page,omitempty" validate:"required"`
	Limit int `query:"limit,omitempty" validate:"required"`
}

// here this are type of repo which are entire seperately handled
type CreateModal struct {
	TCreate
	UserId           string
	VerifiedPurchase bool
	Status           ReviewStatus `json:"status" validate:"required"`
}

type FilterListModel struct {
	Page  int `json:"page,omitempty"`
	Limit int `json:"limit,omitempty"`
}
