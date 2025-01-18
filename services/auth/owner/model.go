package owner

import (
	"time"

	"github.com/omkarp02/pro/services/utils/store"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Owner struct {
	ID          bson.ObjectID    `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string           `json:"name,omitempty" bson:"name"`
	Email       string           `json:"email,omitempty" bson:"email"`
	FirstName   string           `json:"firstname,omitempty" bson:"firstname,omitempty"`
	LastName    string           `json:"lastname,omitempty" bson:"lastname,omitempty"`
	Businesses  []bson.ObjectID  `json:"businesses,omitempty" bson:"businesses"`
	DateOfBirth time.Time        `json:"dateofbirth,omitempty" bson:"dateofbirth,omitempty"`
	Gender      string           `json:"gender,omitempty" bson:"gender,omitempty"`
	MobileNo    string           `json:"mobileNo,omitempty" bson:"mobileNo,omitempty"`
	Timestamps  store.Timestamps `bson:",inline"`
}
