package address

import (
	"context"

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

	result, err := s.txn.RunInTxn(ctx, func(sessCtx context.Context) (interface{}, error) {
		if createPayload.IsPrimary {
			err := s.repo.UpdateAddressIsPrimary(sessCtx, userId)
			if err != nil {
				return "", err
			}
		}

		result, err := s.repo.Create(sessCtx, addressModal)
		if err != nil {
			return "", err
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
	n := 2
	if !payload.IsPrimary {
		n = 1
	}

	go func() {
		respChan <- s.repo.UpdateById(ctx, payload)
	}()

	if payload.IsPrimary {
		go func() {
			respChan <- s.repo.HandleAddressIsPrimary(ctx, payload.UserID, payload.Id)
		}()
	}

	var err error
	for i := 0; i < n; i++ {
		err = <-respChan
	}
	return err
}

func (s *Service) DeleteAddress(ctx context.Context, addressId string, userId string) error {
	return s.repo.Delete(ctx, addressId, userId)
}
func (s *Service) DeleteAddressByIds(ctx context.Context, addressIds []string, userId string) error {
	return s.repo.DeleteByIds(ctx, addressIds, userId)
}

func (s *Service) FindById(ctx context.Context, id string, userId string) (Address, error) {
	project := []string{}

	return s.repo.FindById(ctx, id, userId, project, false)
}

func (s *Service) FindPrimaryAddress(ctx context.Context, userId string) (Address, error) {
	project := []string{}

	return s.repo.FindPrimaryAddress(ctx, userId, project, false)
}
