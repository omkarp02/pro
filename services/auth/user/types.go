package user

import "time"

type TCreateOwner struct {
	Name        string    `json:"name,omitempty"`
	Email       string    `json:"email,omitempty"`
	FirstName   string    `json:"firstname,omitempty"`
	Password    string    `json:"password,omitempty"`
	LastName    string    `json:"lastname,omitempty"`
	DateOfBirth time.Time `json:"dateofbirth,omitempty"`
	MobileNo    string    `json:"mobileNo,omitempty"`
	Gender      string    `json:"gender,omitempty"`
}

type CreateOwnerAndAccountModel struct {
	Name         string    `json:"name,omitempty"`
	Email        string    `json:"email,omitempty"`
	FirstName    string    `json:"firstname,omitempty"`
	Password     string    `json:"password,omitempty"`
	LastName     string    `json:"lastname,omitempty"`
	ProviderId   string    `json:"providerId,omitempty"`
	ProviderName string    `json:"providerName,omitempty"`
	DateOfBirth  time.Time `json:"dateofbirth,omitempty"`
	MobileNo     string    `json:"mobileNo,omitempty"`
	Gender       string    `json:"gender,omitempty"`
}
