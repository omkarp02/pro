package review

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/services/middleware"
	"github.com/omkarp02/pro/utils"
	"github.com/omkarp02/pro/utils/validation"
)

type ReviewService interface {
	Create(ctx context.Context, createBody CreateModal) (string, error)
	Find(ctx context.Context, filterPayload FilterListModel) ([]Review, error)
	FindById(ctx context.Context, id string) (Review, error)
	VoteReview(ctx context.Context, createBody CreateReviewVoteModal) error
}

type Handler struct {
	service   ReviewService
	cfg       *config.Config
	validator *validation.Validator
}

func NewHandler(service ReviewService, cfg *config.Config, validator *validation.Validator) *Handler {
	return &Handler{service: service, cfg: cfg, validator: validator}
}

func (h *Handler) RegisterRoutes(router router.Router, link string) {
	routeGrp := router.Group(link)

	routeGrp.Use(middleware.VerifyToken(h.cfg))

	routeGrp.Post("/", h.create)
	routeGrp.Get("/", h.find)
	routeGrp.Post("/vote", h.VoteReview)
}

func (h *Handler) create(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	var createBody TCreate

	userId := c.GetDecodedData().ID

	if err := h.validator.ValidateBody(c, &createBody); err != nil {
		return err
	}

	modal := CreateModal{
		TCreate:          createBody,
		UserId:           userId,
		VerifiedPurchase: true,
		Status:           REVIEW_STATUS_APPROVED,
	}

	id, err := h.service.Create(ctx, modal)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Address Created Successfully", fiber.Map{"id": id}, 201)
}

func (h *Handler) find(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	var filterData TFilterList

	if err := h.validator.ValidateParams(c, &filterData); err != nil {
		return err
	}
	data, err := h.service.Find(ctx, FilterListModel(filterData))
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Data Fetched Successfully", data, 200)
}

func (h *Handler) VoteReview(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	var body TReviewVoteCreate

	userId := c.GetDecodedData().ID

	if err := h.validator.ValidateBody(c, &body); err != nil {
		return err
	}

	model := CreateReviewVoteModal{
		TReviewVoteCreate: body,
		UserId:            userId,
	}

	err := h.service.VoteReview(ctx, model)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "Data Fetched Successfully", nil, 200)
}

func createContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
