package owner

import (
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/utils/validation"
)

func Intialize(curDb *db.Database, cfg *config.Config, validator *validation.Validator, api router.Router) {

	authConfig := cfg.App.Auth

	// txn := db.NewMongoTransactionManager(curDb)

	repo := NewRepo(curDb, authConfig.DBCollection.Owner)
	service := NewService(repo)
	routeHandler := NewHandler(service, cfg, validator)
	routeHandler.RegisterRoutes(api, authConfig.Routes.Owner)
}
