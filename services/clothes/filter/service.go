package filter

import "context"

type Service struct {
	filterRepo     *FilterRepo
	filterTypeRepo *FilterTypeRepo
}

func NewService(filterRepo *FilterRepo, filterTypeRepo *FilterTypeRepo) *Service {
	return &Service{
		filterRepo:     filterRepo,
		filterTypeRepo: filterTypeRepo,
	}
}

func (s *Service) CreateFilter(ctx context.Context, createFilter TCreateFilter) (string, error) {
	return s.filterRepo.Create(ctx, CreateFilterModal(createFilter))
}

func (s *Service) FindFitlerType(ctx context.Context, filterPayload FilterListModel) ([]FilterType, error) {
	project := []string{}
	return s.filterTypeRepo.FindByFilter(ctx, filterPayload, project, false)
}

func (s *Service) CreateFilterType(ctx context.Context, createFilterType TCreateFilterType) (string, error) {
	return s.filterTypeRepo.Create(ctx, CreateFilterTypeModal(createFilterType))
}
