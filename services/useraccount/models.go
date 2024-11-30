package useraccount

import (
	services "github.com/omkarp02/pro/services/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserAccount struct {
	ID           bson.ObjectID       `json:"_id,omitempty" bson:"_id,omitempty"`
	Email        string              `json:"fullname,omitempty"`
	Password     string              `json:"age,omitempty"`
	Timestamps   services.Timestamps `bson:",inline"`
	RefreshToken string              `json:"refresh_token,omitempty" bson:"refresh_token,omitempty"`
}
