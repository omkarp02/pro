package user

import (
	"context"
	"errors"

	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/services/auth/owner"
	"github.com/omkarp02/pro/services/auth/useraccount"
	"github.com/omkarp02/pro/services/auth/userprofile"
	"github.com/omkarp02/pro/utils/constant"
	"github.com/omkarp02/pro/utils/errutil"
)

type Service struct {
	useraccountRepo *useraccount.Repo
	userprofileRepo *userprofile.Repo
	ownerRepo       *owner.Repo
	txn             db.TransactionManager
}

func NewService(useraccountRepo *useraccount.Repo, userprofileRepo *userprofile.Repo, ownerRepo *owner.Repo, txn db.TransactionManager) *Service {
	return &Service{
		useraccountRepo,
		userprofileRepo,
		ownerRepo,
		txn,
	}
}

func (s *Service) CreateUserProfileAndAccount(ctx context.Context, userprofile userprofile.CreateUserModel, useraccount useraccount.CreateUserAccountModal) (string, error) {

	result, err := s.txn.RunInTxn(ctx, func(sessCtx context.Context) (interface{}, error) {
		userProfileId, err := s.userprofileRepo.Create(ctx, userprofile)
		if err != nil {
			return "", err
		}

		useraccount.UserProfile = userProfileId

		return s.useraccountRepo.Create(ctx, useraccount)
	})

	if err != nil {
		return "", err
	}

	data, _ := result.(string)

	return data, err

}

func (s *Service) CreateUserProfile(ctx context.Context, paylaod userprofile.TCreateUser, useraccountId string) (string, error) {

	result, err := s.txn.RunInTxn(ctx, func(sessCtx context.Context) (interface{}, error) {

		useraccountDetails, err := s.useraccountRepo.FindById(ctx, useraccountId, []string{"email"}, true)
		if err != nil {
			return "", err
		}

		paylaod.Email = useraccountDetails.Email

		id, err := s.userprofileRepo.Create(sessCtx, userprofile.CreateUserModel(paylaod))
		if err != nil {
			return "", err
		}

		if err := s.useraccountRepo.UpdateUserProfileById(ctx, useraccountId, id); err != nil {
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

func (s *Service) CreateOwnerAndAccount(ctx context.Context, ownerPayload CreateOwnerAndAccountModel, userId string) (string, error) {

	result, err := s.txn.RunInTxn(ctx, func(sessCtx context.Context) (interface{}, error) {

		ownerFormattedData := owner.CreateModal{
			TCreateOwner: owner.TCreateOwner{
				Name:        ownerPayload.Name,
				Email:       ownerPayload.Email,
				FirstName:   ownerPayload.FirstName,
				LastName:    ownerPayload.LastName,
				DateOfBirth: ownerPayload.DateOfBirth,
				MobileNo:    ownerPayload.MobileNo,
				Gender:      ownerPayload.Gender,
			},
		}

		createUserAccountModal := useraccount.CreateUserAccountModal{
			Email: ownerPayload.Email,
			AuthProvider: []useraccount.AuthProviderType{
				{
					Provider:   ownerPayload.ProviderName,
					ProviderID: ownerPayload.ProviderName,
				},
			},
			Role: []string{constant.ROLE_OWNER},
		}

		ownerId, err := s.ownerRepo.Create(ctx, ownerFormattedData)
		if err != nil {
			return "", err
		}

		createUserAccountModal.UserProfile = ownerId

		return s.useraccountRepo.Create(ctx, createUserAccountModal)
	})

	if err != nil {
		return "", err
	}

	data, _ := result.(string)

	return data, err

}

func (s *Service) GetUser(ctx context.Context, field string, value string) (useraccount.UserAccount, error) {
	return s.useraccountRepo.FindOne(ctx, field, value, []string{}, false)
}

func (s *Service) HandleRefreshTokenForLogin(ctx context.Context, userId string, refreshToken string, oldRefreshToken string) error {
	var action string

	if len(oldRefreshToken) == 0 {
		action = "push"
	} else {
		if err := s.useraccountRepo.PullUserRefreshToken(ctx, oldRefreshToken); err != nil {
			if errors.Is(err, errutil.ErrDocumentNotFound) {
				action = "reinitialize"
			} else {
				return err
			}
		} else {
			action = "push"
		}

	}

	if err := s.useraccountRepo.UpdateUserRefreshToken(ctx, userId, action, refreshToken); err != nil {
		return err
	}

	return nil
}
