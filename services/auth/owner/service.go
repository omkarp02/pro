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

func (s *Service) Create(ctx context.Context, createOwnerBody CreateModal) (string, error) {
	return s.repo.Create(ctx, createOwnerBody)
}

func (s *Service) Find(ctx context.Context, filterPayload FilterOwnerListModel) ([]Owner, error) {

	project := []string{}

	return s.repo.FindByFilter(ctx, filterPayload, project, false)
}
