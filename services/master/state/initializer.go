package state

import (
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/utils/validation"
)

func Intialize(curDb *db.Database, cfg *config.Config, validator *validation.Validator, api router.Router) {

	configuration := cfg.App.Master

	txn := db.NewMongoTransactionManager(curDb)

	repo := NewRepo(curDb, configuration.DBCollection.State)
	service := NewService(repo, txn)
	routeHandler := NewHandler(service, cfg, validator)
	routeHandler.RegisterRoutes(api, configuration.Routes.Master)
}
