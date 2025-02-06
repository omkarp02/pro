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
		Type:      TYPE_HOME,
		UserID:    userId,
	}

	return s.repo.Create(ctx, addressModal)
}

func (s *Service) GetAddressByUserId(ctx context.Context, userId string) ([]Address, error) {
	return s.repo.GetAddressByUserId(ctx, userId)
}

func (s *Service) UpdateAddress(ctx context.Context, payload UpdateAddressModel) error {
	return s.repo.UpdateById(ctx, payload)
}

func (s *Service) DeleteAddress(ctx context.Context, addressId string, userId string) error {
	return s.repo.Delete(ctx, addressId, userId)
}
