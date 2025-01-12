package userprofile

import "time"

type TCreateUser struct {
	Email       string    `json:"email,omitempty" validate:"required"`
	FirstName   string    `json:"firstname,omitempty" validate:"required"`
	LastName    string    `json:"lastname,omitempty" validate:"required"`
	DateOfBirth time.Time `json:"dateofbirth,omitempty" validate:"required"`
	Gender      string    `json:"gender,omitempty" validate:"required"`
}

type CreateUserModel struct {
	FirstName   string    `json:"firstname,omitempty"`
	LastName    string    `json:"lastname,omitempty"`
	DateOfBirth time.Time `json:"dateofbirth,omitempty"`
	Gender      string    `json:"gender,omitempty"`
	Email       string    `json:"email,omitempty"`
}
