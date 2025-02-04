package review

import (
	"context"
	"fmt"

	"github.com/omkarp02/pro/db"
)

type Service struct {
	repo           *Repo
	reviewVoteRepo *ReviewVoteRepo
	txn            db.TransactionManager
}

func NewService(repo *Repo, reviewVoteRepo *ReviewVoteRepo, txn db.TransactionManager) *Service {
	return &Service{
		repo:           repo,
		reviewVoteRepo: reviewVoteRepo,
		txn:            txn,
	}
}

func (s *Service) Create(ctx context.Context, createBody CreateModal) (string, error) {
	return s.repo.Create(ctx, createBody)
}

func (s *Service) Find(ctx context.Context, filterPayload FilterListModel) ([]Review, error) {
	project := []string{}
	return s.repo.FindByFilter(ctx, filterPayload, project, false)
}

func (s *Service) FindById(ctx context.Context, id string) (Review, error) {
	project := []string{}
	return s.repo.FindById(ctx, id, project, false)
}

func (s *Service) VoteReview(ctx context.Context, createBody CreateReviewVoteModal) error {
	_, err := s.txn.RunInTxn(ctx, func(sessCtx context.Context) (interface{}, error) {

		fmt.Println(createBody, "<<<<<<<<<< here is revieid")
		if err := s.repo.UpdateVotes(sessCtx, createBody.ReviewId, createBody.IsHelpful); err != nil {
			return "", err
		}

		fmt.Println("reached here")

		_, err := s.reviewVoteRepo.Create(sessCtx, createBody)
		return "", err
	})
	return err
}
