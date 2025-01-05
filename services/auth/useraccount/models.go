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
	Email        string           `json:"fullname,omitempty"`
	PasswordHash string           `json:"age,omitempty"`
	AuthProvider []AuthProvider   `bson:"auth_providers" json:"auth_providers"`
	UserProfile  bson.ObjectID    `bson:"userProfileId"`
	Timestamps   store.Timestamps `bson:",inline"`
}
