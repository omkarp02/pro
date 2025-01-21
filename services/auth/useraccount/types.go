package useraccount

type TCreateUserAccount struct {
	Email    string `json:"email" validate:"required_without=PhoneNo"`
	PhoneNo  string `json:"phoneNo" validate:"required_without=Email"`
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
	Email        string
	PasswordHash string
	UserProfile  string
	AuthProvider []AuthProviderType
	Role         []string
	PhoneNumber  string
}

type UpdateUserAccountModel struct {
	UserProfile string
}
