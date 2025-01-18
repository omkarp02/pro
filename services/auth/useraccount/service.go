package useraccount

import (
	"context"
	"errors"
	"fmt"

	"github.com/omkarp02/pro/utils/errutil"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateUser(ctx context.Context, payload CreateUserAccountModal) (string, error) {
	return s.repo.Create(ctx, payload)
}

func (s *Service) HandleRefreshTokenForLogin(ctx context.Context, userId string, refreshToken string, oldRefreshToken string) error {
	var action string

	fmt.Println(userId, refreshToken, oldRefreshToken)

	if len(oldRefreshToken) == 0 {
		action = "push"
	} else {
		if err := s.repo.PullUserRefreshToken(ctx, oldRefreshToken); err != nil {
			if errors.Is(err, errutil.ErrDocumentNotFound) {
				action = "reinitialize"
			} else {
				return err
			}
		} else {
			action = "push"
		}

	}

	if err := s.repo.UpdateUserRefreshToken(ctx, userId, action, refreshToken); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetUser(ctx context.Context, field string, value string) (UserAccount, error) {
	return s.repo.FindOne(ctx, field, value, []string{}, false)
}

func (s *Service) GetUserById(ctx context.Context, id string) (UserAccount, error) {
	return s.repo.FindById(ctx, id, []string{}, false)
}

func (s *Service) UpdateUserRefreshToken(ctx context.Context, userId string, action string, refreshToken string) error {
	return s.repo.UpdateUserRefreshToken(ctx, userId, action, refreshToken)
}

func (s *Service) PullUserRefreshToken(ctx context.Context, refreshToken string) error {
	return s.repo.PullUserRefreshToken(ctx, refreshToken)
}
