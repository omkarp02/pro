package userprofile

import (
	"context"

	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/services/auth/useraccount"
)

type Service struct {
	repo     *Repo
	userRepo *useraccount.Store
	txn      db.TransactionManager
}

func NewService(repo *Repo, userRepo *useraccount.Store, txn db.TransactionManager) *Service {
	return &Service{
		repo:     repo,
		userRepo: userRepo,
		txn:      txn,
	}
}

func (s *Service) CreateUser(ctx context.Context, createUserPayload TCreateUser, userAccountId string) (string, error) {

	emailId, err := s.userRepo.GetUserAccountEmailById(userAccountId)
	if err != nil {
		return "", err
	}

	result, err := s.txn.RunInTxn(ctx, func(sessCtx context.Context) (interface{}, error) {
		payload := CreateUserModel{
			Email:       emailId,
			FirstName:   createUserPayload.FirstName,
			LastName:    createUserPayload.LastName,
			DateOfBirth: createUserPayload.DateOfBirth,
			Gender:      createUserPayload.Gender,
		}

		id, err := s.repo.CreateUser(sessCtx, payload)
		if err != nil {
			return "", err
		}

		err = s.userRepo.UpdateUserAccountProfileById(sessCtx, userAccountId, id)
		if err != nil {
			return "", err
		}

		return id, err
	})

	data, _ := result.(string)

	return data, err
}
