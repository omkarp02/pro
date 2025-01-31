package master

import (
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/utils/validation"
)

func Intialize(curDb *db.Database, cfg *config.Config, validator *validation.Validator, api router.Router) {

	configuration := cfg.App.Master

	repo := NewRepo(curDb, configuration.DBCollection.Master)
	service := NewService(repo)
	routeHandler := NewHandler(service, cfg, validator)
	routeHandler.RegisterRoutes(api, configuration.Routes.Master)
}
