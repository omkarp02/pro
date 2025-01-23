package master

import (
	"context"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, createBody CreateModal) (string, error) {
	return s.repo.Create(ctx, createBody)
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
