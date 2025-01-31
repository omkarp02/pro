package city

import (
	"context"

	"github.com/omkarp02/pro/db"
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

func (s *Service) Create(ctx context.Context, createBody CreateModal) (string, error) {
	return s.repo.Create(ctx, createBody)
}

func (s *Service) Find(ctx context.Context, filterPayload FilterListModel) ([]TCityRes, error) {

	project := []string{"_id", "stateId", "name", "stateDetails"}

	return s.repo.FindByFilter(ctx, filterPayload, project, true)
}

func (s *Service) FindById(ctx context.Context, id string) (City, error) {

	project := []string{}

	return s.repo.FindById(ctx, id, project, false)
}

// objectIds, err := store.SliceOfHexToObjectID(bussinessId, updatedBy)
// 	if err != nil {
// 		return err
// 	}

// result, err := s.txn.RunInTxn(ctx, func(sessCtx context.Context) (interface{}, error) {})

// if err != nil {
// 	return err
// }

// data, _ := result.(string)

// return data, nil
