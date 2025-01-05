package useraccount

type CreateUserAccountBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginUserAccountType struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

//from here this are store type

type AuthProviderType struct {
	Provider     string   `json:"provider"`    // e.g., "google", "facebook"
	ProviderID   string   `json:"provider_id"` // Unique ID from the provider
	RefreshToken []string `json:"refresh_token,omitempty"`
}

type CreateUserAccountModal struct {
	Email        string             `json:"fullname,omitempty"`
	PasswordHash string             `json:"age,omitempty"`
	AuthProvider []AuthProviderType `json:"auth_providers"`
}
