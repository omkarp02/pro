package address

import (
	"context"

	"github.com/omkarp02/pro/services/utils/store"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, userId string, createPayload TCreateAddress) (string, error) {

	addressModal := CreateAddressModel{
		Address:   store.AddressModel(createPayload.Address),
		IsPrimary: createPayload.IsPrimary,
		Type:      createPayload.Type,
		UserID:    userId,
	}

	return s.repo.Create(ctx, addressModal)
}

func (s *Service) GetAddressByUserId(ctx context.Context, userId string) ([]Address, error) {
	return s.repo.GetAddressByUserId(ctx, userId)
}
