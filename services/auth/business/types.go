package bussiness

import "github.com/omkarp02/pro/services/utils/store"

type TCreate struct {
	Name        string         `json:"name,omitempty" validate:"required"`
	Category    string         `json:"category,omitempty" validate:"required"`
	Description string         `json:"description,omitempty" validate:"required"`
	Address     store.TAddress `json:"address,omitempty" validate:"required"`
	Contacts    []Contact      `json:"contacts,omitempty" validate:"required"`
	Website     string         `json:"website,omitempty" validate:"required"`
	LogoUrl     string         `json:"logoUrl,omitempty" validate:"required"`
}

type TFilterList struct {
	Name  string `query:"name,omitempty"`
	Page  int    `query:"page,omitempty" validate:"required"`
	Limit int    `query:"limit,omitempty" validate:"required"`
}

type CreateModal struct {
	Name        string
	OwnerID     string
	Category    string
	Description string
	Address     store.AddressModel
	Contacts    []Contact
	Website     string
	LogoUrl     string
	Active      bool
	CreatorId   string
}

type FilterListModel struct {
	Name  string `json:"name,omitempty"`
	Page  int    `json:"page,omitempty"`
	Limit int    `json:"limit,omitempty"`
}
