package bussiness

import "github.com/omkarp02/pro/services/utils/store"

type CreateBusinessModel struct {
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	OwnerID     string             `json:"ownerId,omitempty" bson:"ownerId,omitempty"`
	Category    string             `json:"category,omitempty" bson:"category,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Address     store.AddressModel `json:"address,omitempty" bson:"address,omitempty"`
	Contacts    []Contact          `json:"contacts,omitempty" bson:"contacts,omitempty"`
	Website     string             `json:"website,omitempty" bson:"website,omitempty"`
	LogoUrl     string             `json:"logoUrl,omitempty" bson:"logoUrl,omitempty"`
	Active      bool               `json:"active,omitempty" bson:"active,omitempty"`
}

// type FilterProductListModel struct {
// 	Sizes      []string `json:"sizes,omitempty"`
// 	Color      string   `json:"color,omitempty"`
// 	MinPrice   float64  `json:"min_price,omitempty"`
// 	MaxPrice   float64  `json:"max_price,omitempty"`
// 	Collection string   `json:"collection,omitempty"`
// 	Name       string   `json:"name,omitempty"`
// 	Page       int      `json:"page,omitempty"`
// 	Limit      int      `json:"limit,omitempty"`
// }
