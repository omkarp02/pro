package store

import "time"

type Timestamps struct {
	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
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
	PinCode           string `bson:"pincode,omitempty"`
	MobileNo          string `bson:"mobileNo,omitempty"`
	AlternateMobileNo string `bson:"alternateMobileNo,omitempty"`
}
