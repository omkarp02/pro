package owner

import "context"

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, createOwnerBody CreateOwnerBody) (string, error) {

	return s.repo.Create(ctx, CreateOwnerModal(createOwnerBody))
}
