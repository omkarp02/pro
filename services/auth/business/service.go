package bussiness

import (
	"context"

	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/services/auth/owner"
)

type Service struct {
	ownerRepo *owner.Repo
	repo      *Repo
	txn       db.TransactionManager
}

func NewService(repo *Repo, ownerRepo *owner.Repo, txn db.TransactionManager) *Service {
	return &Service{
		repo:      repo,
		ownerRepo: ownerRepo,
		txn:       txn,
	}
}

func (s *Service) Create(ctx context.Context, createBody CreateModal, userId string) (string, error) {

	result, err := s.txn.RunInTxn(ctx, func(sessCtx context.Context) (interface{}, error) {

		createBody.Active = false
		businessId, err := s.repo.Create(sessCtx, createBody)
		if err != nil {
			return "", err
		}

		if err := s.ownerRepo.AddBussiness(sessCtx, createBody.OwnerID, userId, businessId); err != nil {
			return "", err
		}

		return businessId, nil
	})

	if err != nil {
		return "", err
	}

	data, _ := result.(string)

	return data, nil
}

func (s *Service) Find(ctx context.Context, filterPayload FilterListModel) ([]Business, error) {

	project := []string{}
	return s.repo.FindByFilter(ctx, filterPayload, project, false)
}

func (s *Service) FindById(ctx context.Context, id string) (Business, error) {

	project := []string{}

	return s.repo.FindById(ctx, id, project, false)
}
