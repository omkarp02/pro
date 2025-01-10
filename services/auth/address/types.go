package address

import "github.com/omkarp02/pro/services/utils/store"

type TCreateAddress struct {
	Address   store.TAddress `json:"address,omitempty"  validate:"required"`
	IsPrimary bool           `json:"is_primary,omitempty"  validate:"required"`
	Type      string         `json:"type,omitempty"  validate:"required"`
}

type CreateAddressModel struct {
	Address   store.AddressModel `json:"address,omitempty"`
	IsPrimary bool               `json:"is_primary,omitempty"`
	Type      string             `json:"type,omitempty"`
	UserID    string             `json:"user_id,omitempty"`
}
