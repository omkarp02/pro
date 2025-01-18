package useraccount

import (
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/utils/validation"
)

func Intialize(curDb *db.Database, cfg *config.Config, validator *validation.Validator, api router.Router) {

	authConfig := cfg.App.Auth

	repo := NewRepo(curDb, authConfig.DBCollection.UserAccount)
	useraccountService := NewService(repo)
	userAccountHandler := NewHandler(useraccountService, cfg, validator)
	userAccountHandler.RegisterRoutes(api, authConfig.Routes.UserAccount)
}
