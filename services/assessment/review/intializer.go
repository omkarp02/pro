package review

import (
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/utils/validation"
)

func Intialize(curDb *db.Database, cfg *config.Config, validator *validation.Validator, api router.Router) {

	configuration := cfg.App.Review

	txn := db.NewMongoTransactionManager(curDb)

	repo := NewRepo(curDb, configuration.DBCollection.Review)
	reviewVoteRepo := NewReviewVoteRepo(curDb, configuration.DBCollection.ReviewVote)
	service := NewService(repo, reviewVoteRepo, txn)
	routeHandler := NewHandler(service, cfg, validator)
	routeHandler.RegisterRoutes(api, configuration.Routes.Review)
}
