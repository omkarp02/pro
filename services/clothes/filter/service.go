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

func (s *Service) CreateFilter(ctx context.Context, createFilter CreateFilterModal) (string, error) {
	return s.filterRepo.Create(ctx, createFilter)
}

func (s *Service) FindFitlerType(ctx context.Context, filterPayload FilterTypeListModel) ([]FilterType, error) {
	project := []string{}
	return s.filterTypeRepo.FindByFilter(ctx, filterPayload, project, false)
}

func (s *Service) FindFitler(ctx context.Context, filterPayload FilterListModel) ([]Filter, error) {
	project := []string{}
	return s.filterRepo.FindByFilter(ctx, filterPayload, project, false)
}

func (s *Service) CreateFilterType(ctx context.Context, createFilterType CreateFilterTypeModal) (string, error) {
	return s.filterTypeRepo.Create(ctx, createFilterType)
}
