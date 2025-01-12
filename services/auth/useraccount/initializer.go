package useraccount

import (
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/services/auth/userprofile"
	"github.com/omkarp02/pro/utils/validation"
)

func Intialize(curDb *db.Database, cfg *config.Config, validator *validation.Validator, api router.Router) {

	authConfig := cfg.App.Auth

	txn := db.NewMongoTransactionManager(curDb)

	userProfileRepo := userprofile.NewRepo(curDb, authConfig.DBCollection.UserProfile)
	userAccountStore := NewStore(curDb, userProfileRepo, authConfig.DBCollection.UserAccount, txn)
	userAccountHandler := NewHandler(userAccountStore, cfg, validator)
	userAccountHandler.RegisterRoutes(api, authConfig.Routes.UserAccount)
}
