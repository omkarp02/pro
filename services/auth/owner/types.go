package owner

import (
	"time"
)

type TCreateOwner struct {
	Name        string    `json:"name,omitempty"`
	Email       string    `json:"email,omitempty"`
	FirstName   string    `json:"firstname,omitempty"`
	LastName    string    `json:"lastname,omitempty"`
	DateOfBirth time.Time `json:"dateofbirth,omitempty"`
	MobileNo    string    `json:"mobileNo,omitempty"`
	Gender      string    `json:"gender,omitempty"`
}

// here this are type of repo which are entire seperately handled
type CreateModal struct {
	TCreateOwner
	CreatorId string
}

type FilterOwnerListModel struct {
	Name  string `json:"name,omitempty"`
	Page  int    `json:"page,omitempty"`
	Limit int    `json:"limit,omitempty"`
}
