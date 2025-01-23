package user

import "time"

type TCreateOwner struct {
	Name        string    `json:"name" validate:"required"`
	Email       string    `json:"email"  validate:"required"`
	FirstName   string    `json:"firstname"  validate:"required"`
	Password    string    `json:"password"  validate:"required"`
	LastName    string    `json:"lastname"  validate:"required"`
	DateOfBirth time.Time `json:"dateofbirth"  validate:"required"`
	MobileNo    string    `json:"mobileNo"  validate:"required"`
	Gender      string    `json:"gender"  validate:"required"`
	Type        string    `json:"type"  validate:"required"`
}

type CreateOwnerAndAccountModel struct {
	Name         string    `json:"name,omitempty"`
	Email        string    `json:"email,omitempty"`
	UserId       string    `json:"userId,omitempty"`
	FirstName    string    `json:"firstname,omitempty"`
	Password     string    `json:"password,omitempty"`
	LastName     string    `json:"lastname,omitempty"`
	ProviderId   string    `json:"providerId,omitempty"`
	ProviderName string    `json:"providerName,omitempty"`
	DateOfBirth  time.Time `json:"dateofbirth,omitempty"`
	MobileNo     string    `json:"mobileNo,omitempty"`
	Gender       string    `json:"gender,omitempty"`
	Type         string
}
