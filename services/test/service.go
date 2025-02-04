package test

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

func (s *Service) Find(ctx context.Context, filterPayload FilterListModel) ([]Model, error) {

	project := []string{}

	return s.repo.FindByFilter(ctx, filterPayload, project, false)
}

func (s *Service) FindById(ctx context.Context, id string) (Model, error) {

	project := []string{}

	return s.repo.FindById(ctx, id, project, false)
}

// func (s *Service) UpdateById(ctx context.Context, id string, userProfileId string) error {

// 	objectId, err := bson.ObjectIDFromHex(userProfileId)
// 	if err != nil {
// 		return err
// 	}

// 	update := bson.M{
// 		"$push": bson.M{
// 			"businesses": objectId, // New name to update
// 		},
// 		"$set": store.GenerateUpdateAudit(objectId),
// 	}
// 	_, err = s.getColl().UpdateByID(ctx, id, update)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

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
