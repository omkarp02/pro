package cart

import (
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/services/clothes/product"
	"github.com/omkarp02/pro/utils/validation"
)

func Intialize(curDb *db.Database, cfg *config.Config, validator *validation.Validator, api router.Router) {

	clothesConfig := cfg.App.Clothes

	CartRepo := NewRepo(curDb, clothesConfig.DBCollection.Cart)
	ProductDetailRepo := product.NewProductDetailRepo(curDb, clothesConfig.DBCollection.ProductDetail)
	CartService := NewService(CartRepo, ProductDetailRepo)
	CartHandler := NewHandler(CartService, cfg, validator)
	CartHandler.RegisterRoutes(api, clothesConfig.Routes.Cart)
}
