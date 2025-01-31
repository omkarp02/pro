package categories

import (
	"context"
	"strconv"

	"github.com/omkarp02/pro/utils/constant"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, createCategory TCreateCategory) (string, error) {

	count, err := s.repo.Count(ctx, FilterCategoryModal{Status: constant.STATUS_ACTIVE})
	if err != nil {
		return "", nil
	}

	modal := CreateCategoryModal{
		CatId:           "CAT-" + strconv.Itoa(count+1),
		TCreateCategory: createCategory,
	}

	return s.repo.Create(ctx, modal)
}

func (s *Service) GetAllCategory(ctx context.Context, filterData FilterCategoryModal, project []string, inclusive bool) ([]TCategoryList, error) {

	var categoryList []TCategoryList

	categoryData, err := s.repo.Find(ctx, filterData, project, inclusive)
	if err != nil {
		return nil, err
	}
	for _, item := range categoryData {

		categoryList = append(categoryList, TCategoryList{
			CatId: item.CatId,
			Name:  item.Name,
			Icon:  item.Icon,
			Slug:  item.Slug,
			Id:    item.ID.Hex(),
		})
	}
	return categoryList, nil
}
