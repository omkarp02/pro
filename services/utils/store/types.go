package store

type AddressModel struct {
	Name              string `bson:"name,omitempty"`
	Address           string `json:"address,omitempty"`
	City              string `json:"city,omitempty"`
	State             string `json:"state,omitempty"`
	Country           string `json:"country,omitempty"`
	PinCode           int    `json:"pincode,omitempty"`
	MobileNo          string `json:"mobileNo,omitempty"`
	AlternateMobileNo string `json:"alternateMobileNo,omitempty"`
}

type TAddress struct {
	Name              string `bson:"name,omitempty" validate:"required"`
	Address           string `json:"address,omitempty"  validate:"required"`
	City              string `json:"city,omitempty"  validate:"required"`
	State             string `json:"state,omitempty"  validate:"required"`
	Country           string `json:"country,omitempty"  validate:"required"`
	PinCode           int    `json:"pincode,omitempty"  validate:"required"`
	MobileNo          string `json:"mobileNo,omitempty"  validate:"required"`
	AlternateMobileNo string `json:"alternateMobileNo,omitempty"`
}

// {
// 	"address": "13 Main street near brooklyn",
// 	"city": "Hong Kong",
// 	"state": "Bancong",
// 	"country": "some",
// 	"pincode": "433534",
// 	"mobileNo": "2348334343",
// 	"alternateMobileNo": "7348334343",
// }
