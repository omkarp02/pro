package userprofile

import (
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/services/auth/useraccount"
	"github.com/omkarp02/pro/utils/validation"
)

func Intialize(curDb *db.Database, cfg *config.Config, validator *validation.Validator, api router.Router) {

	authConfig := cfg.App.Auth

	txn := db.NewMongoTransactionManager(curDb)

	userProfileRepo := NewRepo(curDb, authConfig.DBCollection.UserProfile)
	userAccountStore := useraccount.NewStore(curDb, authConfig.DBCollection.UserAccount)
	userProfileService := NewService(userProfileRepo, userAccountStore, txn)
	userProfileHandler := NewHandler(userProfileService, cfg, validator)
	userProfileHandler.RegisterRoutes(api, authConfig.Routes.UserProfile)
}
