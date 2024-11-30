package services

import "time"

type Timestamps struct {
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type ACCESS_TOKEN_PAYLOAD struct {
	ID string
}

type REFRESH_TOKEN_PAYLOAD struct {
	ID string
}
