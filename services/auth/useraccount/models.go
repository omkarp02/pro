package useraccount

import (
	"github.com/omkarp02/pro/services/utils/store"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type AuthProvider struct {
	Provider   string `bson:"provider" json:"provider"`       // e.g., "google", "facebook"
	ProviderID string `bson:"provider_id" json:"provider_id"` // Unique ID from the provider
}

type UserAccount struct {
	ID           bson.ObjectID    `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId       string           `json:"userId,omitempty" bson:"userId,omitempty"`
	PasswordHash string           `json:"password,omitempty" bson:"password,omitempty"`
	Role         []string         `json:"role,omitempty" bson:"role,omitempty"`
	Type         string           `bson:"type" json:"type"`
	AuthProvider []AuthProvider   `bson:"auth_providers" json:"auth_providers"`
	UserProfile  *bson.ObjectID   `bson:"userProfileId,omitempty" json:"userProfileId"`
	RefreshToken []string         `bson:"refresh_token,omitempty"`
	Timestamps   store.Timestamps `bson:",inline"`
}
