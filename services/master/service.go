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

func (s *Service) Find(ctx context.Context, filterPayload FilterListModel) ([]Master, error) {

	project := []string{"name"}

	return s.repo.FindByFilter(ctx, filterPayload, project, false)
}
