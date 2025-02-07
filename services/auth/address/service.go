package address

import (
	"context"
	"fmt"

	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/services/utils/store"
)

type Service struct {
	repo *Repo
	txn  db.TransactionManager
}

func NewService(repo *Repo, txn db.TransactionManager) *Service {
	return &Service{
		repo: repo,
		txn:  txn,
	}
}

func (s *Service) Create(ctx context.Context, userId string, createPayload TCreateAddress) (string, error) {

	addressModal := CreateAddressModel{
		Address:   store.AddressModel(createPayload.Address),
		IsPrimary: createPayload.IsPrimary,
		Type:      TYPE_HOME,
		UserID:    userId,
	}

	fmt.Println(addressModal, "<<<<<<<<")

	result, err := s.txn.RunInTxn(ctx, func(sessCtx context.Context) (interface{}, error) {
		result, err := s.repo.Create(sessCtx, addressModal)
		if err != nil {
			return "", err
		}

		if createPayload.IsPrimary {
			err = s.repo.UpdateAddressIsPrimary(sessCtx, userId)
		}

		return result, err
	})

	if err != nil {
		return "", err
	}

	data, _ := result.(string)

	return data, nil

}

func (s *Service) GetAddressByUserId(ctx context.Context, userId string) ([]Address, error) {
	return s.repo.GetAddressByUserId(ctx, userId)
}

func (s *Service) UpdateAddress(ctx context.Context, payload UpdateAddressModel) error {

	respChan := make(chan error)

	go func() {
		respChan <- s.repo.UpdateById(ctx, payload)
	}()

	if payload.IsPrimary {
		go func() {
			respChan <- s.repo.HandleAddressIsPrimary(ctx, payload.UserID, payload.Id)
		}()
	}

	var err error
	for i := 0; i < 2; i++ {
		err = <-respChan
	}
	return err
}

func (s *Service) DeleteAddress(ctx context.Context, addressId string, userId string) error {
	return s.repo.Delete(ctx, addressId, userId)
}
