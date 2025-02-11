package filter

import (
	"context"
	"fmt"
)

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

func (s *Service) FindFitlerForUser(ctx context.Context, filterPayload FilterListModel) ([]FilterListForUserRes, error) {

	var finalRes []FilterListForUserRes

	project := []string{}
	filterList, err := s.filterRepo.FindByFilter(ctx, filterPayload, project, false)
	if err != nil {
		return nil, err
	}

	fmt.Println(filterList, "<<<<<<<<<<<<<<<<<")

	filterTypeList, err := s.filterTypeRepo.FindByFilter(ctx, FilterTypeListModel{}, project, false)
	if err != nil {
		return nil, err
	}

	m := make(map[string][]FilterItemForUserRes, len(filterList))

	for _, filterItem := range filterList {
		if _, ok := m[filterItem.Type.Hex()]; ok {
			m[filterItem.Type.Hex()] = append(m[filterItem.Type.Hex()], FilterItemForUserRes{
				Name: filterItem.Name,
				ID:   filterItem.Slug,
			})
		} else {
			m[filterItem.Type.Hex()] = []FilterItemForUserRes{
				{
					Name: filterItem.Name,
					ID:   filterItem.Slug,
				},
			}
		}

	}

	for _, filterTypeItem := range filterTypeList {

		if value, ok := m[filterTypeItem.ID.Hex()]; ok {
			finalRes = append(finalRes, FilterListForUserRes{
				Name:   filterTypeItem.Name,
				Fitler: value,
			})
		}

	}

	return finalRes, nil

}

func (s *Service) CreateFilterType(ctx context.Context, createFilterType CreateFilterTypeModal) (string, error) {
	return s.filterTypeRepo.Create(ctx, createFilterType)
}
