package test

import "time"

type TCreate struct {
	Name        string    `json:"name,omitempty"`
	Email       string    `json:"email,omitempty"`
	FirstName   string    `json:"firstname,omitempty"`
	LastName    string    `json:"lastname,omitempty"`
	DateOfBirth time.Time `json:"dateofbirth,omitempty"`
	MobileNo    string    `json:"mobileNo,omitempty"`
	Gender      string    `json:"gender,omitempty"`
}

type TFilterList struct {
	Name  string `query:"name,omitempty"`
	Page  int    `query:"page,omitempty" validate:"required"`
	Limit int    `query:"limit,omitempty" validate:"required"`
}

// here this are type of repo which are entire seperately handled
type CreateModal struct {
	Name        string    `json:"name,omitempty"`
	Email       string    `json:"email,omitempty"`
	FirstName   string    `json:"firstname,omitempty"`
	LastName    string    `json:"lastname,omitempty"`
	DateOfBirth time.Time `json:"dateofbirth,omitempty"`
	MobileNo    string    `json:"mobileNo,omitempty"`
	Gender      string    `json:"gender,omitempty"`
}

type FilterListModel struct {
	Name  string `json:"name,omitempty"`
	Page  int    `json:"page,omitempty"`
	Limit int    `json:"limit,omitempty"`
}
