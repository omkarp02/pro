package product

import (
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/utils/validation"
)

func Intialize(curDb *db.Database, cfg *config.Config, validator *validation.Validator, api router.Router) {

	productConfig := cfg.App.Clothes

	txn := db.NewMongoTransactionManager(curDb)

	productListRepo := NewProductListRepo(curDb, productConfig.DBCollection.ProductList)
	productDetailRepo := NewProductDetailRepo(curDb, productConfig.DBCollection.ProductDetail)
	productBatchRepo := NewProductBatchRepo(curDb, productConfig.DBCollection.ProductBatch)
	productService := NewService(productListRepo, productDetailRepo, productBatchRepo, txn)
	productHandler := NewHandler(productService, cfg, validator)
	productHandler.RegisterRoutes(api, productConfig.Routes.Product)
}
