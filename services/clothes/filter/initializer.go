package filter

import (
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/utils/validation"
)

func Intialize(curDb *db.Database, cfg *config.Config, validator *validation.Validator, api router.Router) {

	fitlerConfig := cfg.App.Clothes

	filterRepo := NewRepoFilter(curDb, fitlerConfig.DBCollection.Filter)
	filterTypeRepo := NewRepoFilterType(curDb, fitlerConfig.DBCollection.FilterType)
	filterService := NewService(filterRepo, filterTypeRepo)
	filterHandler := NewHandler(filterService, cfg, validator)
	filterHandler.RegisterRoutes(api, fitlerConfig.Routes.Filter)
}
