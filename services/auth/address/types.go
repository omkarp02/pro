package address

import "github.com/omkarp02/pro/services/utils/store"

type TCreateAddress struct {
	Address   store.TAddress `json:"address,omitempty"  validate:"required"`
	Type      ADDRESS_TYPE   `json:"type,omitempty" validate:"required"`
	IsPrimary bool           `json:"isPrimary,omitempty"  validate:"required"`
}

type CreateAddressModel struct {
	Address   store.AddressModel `json:"address,omitempty"`
	IsPrimary bool               `json:"isPrimary,omitempty"`
	Type      ADDRESS_TYPE       `json:"type,omitempty"`
	UserID    string             `json:"user_id,omitempty"`
}

type TUpdateAddress struct {
	Id        string         `json:"id,omitempty"  validate:"required"`
	Address   store.TAddress `json:"address,omitempty"  validate:"required"`
	IsPrimary bool           `json:"isPrimary,omitempty"  validate:"required"`
}

type UpdateAddressModel struct {
	Id        string
	Address   store.AddressModel
	IsPrimary bool
	UserID    string
}
