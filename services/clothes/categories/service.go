package categories

import (
	"context"
	"fmt"
	"strconv"
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
	fmt.Println("reached sercie<<<<<<<<<", createCategory)

	count, err := s.repo.Count(ctx, FilterCategoryModal{IsActive: true})
	if err != nil {
		return "", nil
	}

	modal := CreateCategoryModal{
		CatId:           "CAT-" + strconv.Itoa(count+1),
		TCreateCategory: createCategory,
	}

	return s.repo.Create(ctx, modal)
}

func (s *Service) GetAllCategory(ctx context.Context, filterData TFilterCategory) ([]TCategoryList, error) {

	var categoryList []TCategoryList

	project := []string{"catId", "icon", "name", "slug"}

	categoryData, err := s.repo.Find(ctx, FilterCategoryModal{IsActive: true}, project, true)
	if err != nil {
		return nil, err
	}

	for _, item := range categoryData {

		categoryList = append(categoryList, TCategoryList{
			CatId: item.CatId,
			Name:  item.Name,
			Icon:  item.Icon,
			Slug:  item.Slug,
		})
	}

	return categoryList, nil
}
