package owner

type CreateOwnerBody struct {
	Name       string   `json:"name" validate:"required"`
	Email      string   `json:"email" validate:"required,email"`
	Businesses []string `json:"businesses" validate:"required"`
}

// here this are type of repo which are entire seperately handled
type CreateOwnerModal struct {
	Name       string   `json:"name,omitempty"`
	Email      string   `json:"email,omitempty"`
	Businesses []string `json:"businesses,omitempty"`
}
