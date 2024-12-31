package product

import (
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/utils/validation"
)

func Intialize(curDb *db.Database, collNameProductList, collNameProductDetail string, cfg *config.Config, validator *validation.Validator, routeName string, api router.Router) {
	productListRepo := NewProductListRepo(curDb, collNameProductList)
	productDetailRepo := NewProductDetailRepo(curDb, collNameProductDetail)
	productService := NewService(productListRepo, productDetailRepo)
	productHandler := NewHandler(productService, cfg, validator)
	productHandler.RegisterRoutes(api, routeName)
}
