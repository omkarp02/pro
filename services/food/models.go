package food

import (
	"time"

	"github.com/omkarp02/pro/services/utils/store"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Size struct {
	Size          string  `json:"size,omitempty" bson:"size,omitempty"`
	PriceModifier float64 `json:"price_modifier,omitempty" bson:"price_modifier,omitempty"`
}

type Volume struct {
	Amount float64 `json:"amount,omitempty" bson:"amount,omitempty"`
	Unit   string  `json:"unit,omitempty" bson:"unit,omitempty"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude,omitempty" bson:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty" bson:"longitude,omitempty"`
}

type Location struct {
	RestaurantID string      `json:"restaurant_id,omitempty" bson:"restaurant_id,omitempty"`
	Address      string      `json:"address,omitempty" bson:"address,omitempty"`
	Coordinates  Coordinates `json:"coordinates,omitempty" bson:"coordinates,omitempty"`
}

type NutritionalInfo struct {
	Calories float64 `json:"calories,omitempty" bson:"calories,omitempty"`
	Proteins float64 `json:"proteins,omitempty" bson:"proteins,omitempty"`
	Carbs    float64 `json:"carbs,omitempty" bson:"carbs,omitempty"`
	Fats     float64 `json:"fats,omitempty" bson:"fats,omitempty"`
}

type Image struct {
	URL     string `json:"url,omitempty" bson:"url,omitempty"`
	AltText string `json:"alt_text,omitempty" bson:"alt_text,omitempty"`
}

type CustomizationOption struct {
	Option        string  `json:"option,omitempty" bson:"option,omitempty"`
	PriceModifier float64 `json:"price_modifier,omitempty" bson:"price_modifier,omitempty"`
}

type Discount struct {
	Percentage float64   `json:"percentage,omitempty" bson:"percentage,omitempty"`
	ValidTill  time.Time `json:"valid_till,omitempty" bson:"valid_till,omitempty"`
}

type Availability struct {
	IsAvailable bool `json:"is_available,omitempty" bson:"is_available,omitempty"`
	Stock       int  `json:"stock,omitempty" bson:"stock,omitempty"`
}

type Rating struct {
	Average float64 `json:"average,omitempty" bson:"average,omitempty"`
	Count   int     `json:"count,omitempty" bson:"count,omitempty"`
}

type Product struct {
	ID          bson.ObjectID    `json:"_id,omitempty" bson:"_id,omitempty"`
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

// {
// 	"id": "string", // Unique identifier for the product
// 	"name": "string", // Name of the product
// 	"type": "string", // Type of food (e.g., Chinese, Pizza, Burger, Dessert)
// 	"category": "string", // Category within the type (e.g., Veg, Non-Veg, Vegan)
// 	"price": {
// 	  "amount": "number", // Price of the item
// 	},
// 	"sizes": [
// 	  {
// 		"size": "string", // e.g., Small, Medium, Large
// 		"price_modifier": "number" // Additional cost or discount for the size
// 	  }
// 	],
// 	"ingredients": ["string"], // List of key ingredients
// 	"allergens": ["string"], // List of allergens (e.g., nuts, gluten, dairy)
// 	"nutritional_info": {
// 	  "calories": "number",
// 	  "proteins": "number",
// 	  "carbs": "number",
// 	  "fats": "number"
// 	},
// 	"volume": {
// 	  "amount": "number", // e.g., 500
// 	  "unit": "string" // e.g., ml, g, kg
// 	},
// 	"description": "string", // Detailed description of the product
// 	"location": {
// 	  "restaurant_id": "string", // Link to the restaurant offering this product
// 	  "address": "string", // Restaurant address
// 	  "coordinates": {
// 		"latitude": "number",
// 		"longitude": "number"
// 	  }
// 	},
// 	"availability": {
// 	  "is_available": "boolean",
// 	  "stock": "number" // Units available
// 	},
// 	"tags": ["string"], // e.g., Spicy, Popular, New
// 	"images": [
// 	  {
// 		"url": "string", // URL to the image
// 		"alt_text": "string" // Alt text for accessibility
// 	  }
// 	],
// 	"rating": {
// 	  "average": "number", // Average rating
// 	  "count": "number" // Number of reviews
// 	},
// 	"customization_options": [
// 	  {
// 		"option": "string", // e.g., Extra Cheese, No Onions
// 		"price_modifier": "number" // Additional cost or discount
// 	  }
// 	],
// 	"discount": {
// 	  "percentage": "number", // e.g., 10% off
// 	  "valid_till": "string" // Expiration date of the discount
// 	},
// 	"created_at": "string", // ISO 8601 timestamp
// 	"updated_at": "string" // ISO 8601 timestamp
//   }
