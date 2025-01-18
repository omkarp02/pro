package userprofile

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

func (s *Service) CreateUser(ctx context.Context, createUserPayload TCreateUser, userAccountId string) (string, error) {

	result, err := s.txn.RunInTxn(ctx, func(sessCtx context.Context) (interface{}, error) {
		payload := CreateUserModel{
			Email:       createUserPayload.Email,
			FirstName:   createUserPayload.FirstName,
			LastName:    createUserPayload.LastName,
			DateOfBirth: createUserPayload.DateOfBirth,
			Gender:      createUserPayload.Gender,
		}

		id, err := s.repo.Create(sessCtx, payload)
		if err != nil {
			return "", err
		}

		return id, err
	})

	if err != nil {
		return "", err
	}

	data, _ := result.(string)

	return data, err
}
