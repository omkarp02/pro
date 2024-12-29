package food

import "github.com/omkarp02/pro/services/utils/store"

type CreateProduct struct {
	Name        string           `json:"name,omitempty" bson:"name,omitempty"`
	Type        string           `json:"type,omitempty" bson:"type,omitempty"`
	Category    string           `json:"category,omitempty" bson:"category,omitempty"`
	Price       float64          `json:"price,omitempty" bson:"price,omitempty"`
	Sizes       []Size           `json:"sizes,omitempty" bson:"sizes,omitempty"`
	Volume      Volume           `json:"volume,omitempty" bson:"volume,omitempty"`
	Description string           `json:"description,omitempty" bson:"description,omitempty"`
	Location    Location         `json:"location,omitempty" bson:"location,omitempty"`
	Tags        []string         `json:"tags,omitempty" bson:"tags,omitempty"`
	Images      []Image          `json:"images,omitempty" bson:"images,omitempty"`
	Rating      Rating           `json:"rating,omitempty" bson:"rating,omitempty"`
	Timestamps  store.Timestamps `bson:",inline"`
}
