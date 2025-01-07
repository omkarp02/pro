package useraccount

import (
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/utils/validation"
)

func Intialize(curDb *db.Database, cfg *config.Config, validator *validation.Validator, api router.Router) {

	authConfig := cfg.App.Auth

	userAccountStore := NewStore(curDb, authConfig.DBCollection.UserAccount)
	userAccountHandler := NewHandler(userAccountStore, cfg, validator)
	userAccountHandler.RegisterRoutes(api, authConfig.Routes.UserAccount)
}
