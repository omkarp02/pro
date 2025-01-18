package user

import (
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/services/auth/owner"
	"github.com/omkarp02/pro/services/auth/useraccount"
	"github.com/omkarp02/pro/services/auth/userprofile"
	"github.com/omkarp02/pro/utils/validation"
)

func Intialize(curDb *db.Database, cfg *config.Config, validator *validation.Validator, api router.Router) {

	authConfig := cfg.App.Auth

	txn := db.NewMongoTransactionManager(curDb)

	useraccountRepo := useraccount.NewRepo(curDb, authConfig.DBCollection.UserAccount)
	userprofileRepo := userprofile.NewRepo(curDb, authConfig.DBCollection.UserProfile)
	ownerRepo := owner.NewRepo(curDb, authConfig.DBCollection.Owner)
	userService := NewService(useraccountRepo, userprofileRepo, ownerRepo, txn)
	handler := NewHandler(userService, cfg, validator)
	handler.RegisterRoutes(api, authConfig.Routes.UserAccount)
}
