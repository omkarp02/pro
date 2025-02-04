package review

import (
	"time"

	"github.com/omkarp02/pro/services/utils/store"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Review struct {
	ID               bson.ObjectID     `json:"_id,omitempty" bson:"_id,omitempty"`
	Timestamps       *store.Timestamps `json:"timestamp,omitempty" bson:",inline"`
	UserId           bson.ObjectID     `json:"userId,omitempty" bson:"userId,omitempty"`
	ProductCode      string            `json:"productCode" bson:"productCode"`
	Rating           int               `json:"rating" bson:"rating"`
	ReviewText       string            `json:"reviewText" bson:"reviewText"`
	VerifiedPurchase bool              `json:"verifiedPurchase" bson:"verifiedPurchase"`
	HelpfulVotes     int               `json:"helpfulVotes" bson:"helpfulVotes"`
	NotHelpfulVotes  int               `json:"notHelpfulVotes" bson:"notHelpfulVotes"`
	Status           ReviewStatus      `json:"status" bson:"status"`
}

type ReviewStatus string

const (
	REVIEW_STATUS_APPROVED ReviewStatus = "approved"
	REVIEW_STATUS_PENDING  ReviewStatus = "pending"
	REVIEW_STATUS_FLAGGED  ReviewStatus = "flagged"
)

type ReviewVote struct {
	ID        bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId    bson.ObjectID `json:"userId" bson:"userId"`
	ReviewId  bson.ObjectID `json:"reviewId" bson:"reviewId"`
	IsHelpful bool          `json:"isHelpful" bson:"isHelpful"`
	CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
}
