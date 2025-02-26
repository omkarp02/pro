package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/services/assessment/review"
	"github.com/omkarp02/pro/services/auth/useraccount"
	"github.com/omkarp02/pro/services/clothes/cart"
	"github.com/omkarp02/pro/services/clothes/categories"
	"github.com/omkarp02/pro/services/clothes/filter"
	"github.com/omkarp02/pro/services/clothes/product"
	"github.com/omkarp02/pro/services/master/city"
	"github.com/omkarp02/pro/services/master/master"
	"github.com/omkarp02/pro/services/master/state"
	"github.com/omkarp02/pro/utils/validation"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func main() {

	cfg := config.MustLoad("")
	config.SetUpLogger()

	DB, err := db.NewDatabase(cfg)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err = DB.DB.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	api := router.NewFiberRouter(cfg)
	api.Get("/", func(ctx router.Context) error {
		response := Response{
			Status:  -1,
			Message: "User registered successfully",
			Data:    fiber.Map{"id": "sdflgkj"},
		}

		return ctx.JSON(200, response)
	})

	api.Post("/", func(ctx router.Context) error {
		response := Response{
			Status:  -1,
			Message: "User registered successfully",
			Data:    fiber.Map{"id": "sdflgkj"},
		}

		return ctx.JSON(200, response)
	})

	validator := validation.NewValidator()

	// setUpClothesApp(DB, cfg, validator, api)
	// setUpAuthApp(DB, cfg, validator, api)
	useraccount.Intialize(DB, cfg, validator, api)

	// setUpMasterApp(DB, cfg, validator, api)
	// setUpAssessment(DB, cfg, validator, api)

	api.Listen(cfg.HTTPServer.Addr)
}

func setUpClothesApp(curDb *db.Database, cfg *config.Config, validator *validation.Validator, api router.Router) {
	categories.Intialize(curDb, cfg, validator, api)
	filter.Intialize(curDb, cfg, validator, api)
	product.Intialize(curDb, cfg, validator, api)
	cart.Intialize(curDb, cfg, validator, api)
}

func setUpAuthApp(curDb *db.Database, cfg *config.Config, validator *validation.Validator, api router.Router) {
	useraccount.Intialize(curDb, cfg, validator, api)
	// userprofile.Intialize(curDb, cfg, validator, api)
	// user.Intialize(curDb, cfg, validator, api)
	// owner.Intialize(curDb, cfg, validator, api)
	// bussiness.Intialize(curDb, cfg, validator, api)
	// address.Intialize(curDb, cfg, validator, api)
}

func setUpMasterApp(curDb *db.Database, cfg *config.Config, validator *validation.Validator, api router.Router) {
	master.Intialize(curDb, cfg, validator, api)
	state.Intialize(curDb, cfg, validator, api)
	city.Intialize(curDb, cfg, validator, api)
}

func setUpAssessment(curDb *db.Database, cfg *config.Config, validator *validation.Validator, api router.Router) {
	review.Intialize(curDb, cfg, validator, api)
}
