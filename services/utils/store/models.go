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
	Name              string `json:"name,omitempty" bson:"name,omitempty"`
	Address           string `json:"address,omitempty" bson:"address,omitempty"`
	City              string `json:"city,omitempty" bson:"city,omitempty"`
	State             string `json:"state,omitempty" bson:"state,omitempty"`
	Country           string `json:"country,omitempty" bson:"country,omitempty"`
	PinCode           int    `json:"pincode,omitempty" bson:"pincode,omitempty"`
	MobileNo          string `json:"mobileNo,omitempty" bson:"mobileNo,omitempty"`
	AlternateMobileNo string `json:"alternateMobileNo,omitempty" bson:"alternateMobileNo,omitempty"`
}
