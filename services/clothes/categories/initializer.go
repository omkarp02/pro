package categories

import (
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/utils/validation"
)

func Intialize(curDb *db.Database, cfg *config.Config, validator *validation.Validator, api router.Router) {

	clothesConfig := cfg.App.Clothes

	categoryRepo := NewRepo(curDb, clothesConfig.DBCollection.Category)
	categoryService := NewService(categoryRepo)
	categoryHandler := NewHandler(categoryService, cfg, validator)
	categoryHandler.RegisterRoutes(api, clothesConfig.Routes.Category)
}
