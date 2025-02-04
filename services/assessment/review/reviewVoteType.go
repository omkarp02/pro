package review

type TReviewVoteCreate struct {
	ReviewId  string `json:"reviewId" validate:"required"`
	IsHelpful bool   `json:"isHelpful" validate:"required"`
}

// here this are type of repo which are entire seperately handled
type CreateReviewVoteModal struct {
	TReviewVoteCreate
	UserId string `json:"userId" bson:"userId"`
}

type FilterListReviewVoteModel struct {
	ReviewId string `json:"reviewId" bson:"reviewId"`
}
