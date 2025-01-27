package store

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Timestamps struct {
	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

type AuditFields struct {
	CreatedAt  time.Time     `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	CreatedBy  bson.ObjectID `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
	UpdatedAt  time.Time     `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	ModifiedBy bson.ObjectID `json:"modifiedBy,omitempty" bson:"modifiedBy,omitempty"`
}

type Pagination struct {
	Page  int `json:"page,omitempty"`
	Limit int `json:"limit,omitempty"`
}

type Address struct {
	Address           string `bson:"address,omitempty"`
	City              string `bson:"city,omitempty"`
	State             string `bson:"state,omitempty"`
	Country           string `bson:"country,omitempty"`
	PinCode           int    `bson:"pincode,omitempty"`
	MobileNo          string `bson:"mobileNo,omitempty"`
	AlternateMobileNo string `bson:"alternateMobileNo,omitempty"`
}
