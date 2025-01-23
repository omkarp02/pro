package useraccount

type TCreateUserAccount struct {
	UserId   string `json:"userId"  validate:"required"`
	Type     string `json:"type"  validate:"required,oneof=phone email"`
	Password string `json:"password" validate:"required"`
}

type LoginUserAccountType struct {
	UserId   string `json:"userId" validate:"required"`
	Password string `json:"password" validate:"required"`
}

//from here this are store type

type AuthProviderType struct {
	Provider     string   `json:"provider"`    // e.g., "google", "facebook"
	ProviderID   string   `json:"provider_id"` // Unique ID from the provider
	RefreshToken []string `json:"refresh_token,omitempty"`
}

type CreateUserAccountModal struct {
	UserId       string
	PasswordHash string
	UserProfile  string
	AuthProvider []AuthProviderType
	Role         []string
	Type         string
	PhoneNumber  string
}

type UpdateUserAccountModel struct {
	UserProfile string
}
