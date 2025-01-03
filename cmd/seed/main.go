package main

import (
	"context"
	"log"
	"time"

	"github.com/omkarp02/pro/config"
	data "github.com/omkarp02/pro/data/clothes"
	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/services/clothes/categories"
	"github.com/omkarp02/pro/services/clothes/product"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	var configPath string = ""

	cfg := config.MustLoad(configPath)
	config.SetUpLogger()

	DB, err := db.NewDatabase(cfg)
	if err != nil {
		log.Fatal(err)
	}

	populateCategories(DB, cfg)
	// populateProducts(DB, cfg)

}

func populateCategories(db *db.Database, cfg *config.Config) {
	ctx, cancel := createContext()
	defer cancel()

	for _, cat := range data.CategoriesData {

		category := categories.Category{
			CatId:       cat.Id,
			Name:        cat.Name,
			Description: cat.Description,
			ImgLink:     cat.ImgLink,
			Icon:        cat.Icon,
			IsActive:    cat.IsActive,
			Slug:        cat.Slug,
		}

		update := bson.M{
			"$set": category,
		}

		opts := options.Update().SetUpsert(true)
		filter := bson.M{"catId": category.CatId}

		_, err := db.DB.Database(db.DBName).Collection(cfg.App.Clothes.DBCollection.Category).UpdateOne(ctx, filter, update, opts)
		if err != nil {
			panic(err)
		}
	}
}

func populateProducts(db *db.Database, cfg *config.Config) {
	ctx, cancel := createContext()
	defer cancel()

	productListRepo := product.NewProductListRepo(db, cfg.App.Clothes.DBCollection.ProductList)
	productDetailRepo := product.NewProductDetailRepo(db, cfg.App.Clothes.DBCollection.ProductDetail)
	productService := product.NewService(productListRepo, productDetailRepo)

	for _, p := range data.ProductListData {

		_, err := productService.CreateProductList(ctx, product.TCreateProductList(p))

		if err != nil {
			panic(err)
		}
	}
}

func createContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
