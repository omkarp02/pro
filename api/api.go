package api

import (
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/router"
	bussiness "github.com/omkarp02/pro/services/auth/business"
	"github.com/omkarp02/pro/services/auth/owner"
	"github.com/omkarp02/pro/services/auth/user"
	"github.com/omkarp02/pro/services/auth/useraccount"
	"github.com/omkarp02/pro/services/auth/userprofile"
	"github.com/omkarp02/pro/services/clothes/cart"
	"github.com/omkarp02/pro/services/clothes/categories"
	"github.com/omkarp02/pro/services/clothes/filter"
	"github.com/omkarp02/pro/services/clothes/product"
	"github.com/omkarp02/pro/services/master"
	"github.com/omkarp02/pro/utils/validation"
)

type APIServer struct {
	addr   string
	db     *db.Database
	config *config.Config
}

func NewAPIServer(addr string, curDb *db.Database, config *config.Config) *APIServer {
	return &APIServer{
		addr:   addr,
		db:     curDb,
		config: config,
	}
}

func (s *APIServer) Run() error {

	api := router.NewFiberRouter(s.config)

	validator := validation.NewValidator()

	setUpClothesApp(s.db, s.config, validator, api)
	setUpAuthApp(s.db, s.config, validator, api)
	setUpMasterApp(s.db, s.config, validator, api)

	return api.Listen(s.addr)
}

func setUpClothesApp(curDb *db.Database, cfg *config.Config, validator *validation.Validator, api router.Router) {
	categories.Intialize(curDb, cfg, validator, api)
	filter.Intialize(curDb, cfg, validator, api)
	product.Intialize(curDb, cfg, validator, api)
	cart.Intialize(curDb, cfg, validator, api)
}

func setUpAuthApp(curDb *db.Database, cfg *config.Config, validator *validation.Validator, api router.Router) {
	useraccount.Intialize(curDb, cfg, validator, api)
	userprofile.Intialize(curDb, cfg, validator, api)
	user.Intialize(curDb, cfg, validator, api)
	owner.Intialize(curDb, cfg, validator, api)
	bussiness.Intialize(curDb, cfg, validator, api)
}

func setUpMasterApp(curDb *db.Database, cfg *config.Config, validator *validation.Validator, api router.Router) {
	master.Intialize(curDb, cfg, validator, api)
}

// authHandler := auth.NewHandler(s.config, userAccountStore)
// authHandler.RegisterRoutes(api, "auth")

// ownerRepo := owner.NewRepo(s.db, "owners")
// ownerService := owner.NewService(ownerRepo)
// ownerHandler := owner.NewHandler(ownerService, s.config, validator)
// ownerHandler.RegisterRoutes(api, "owner")

// userProfileStore := userprofile.NewStore(s.db, "user_profile")
// userHandler := userprofile.NewHandler(userProfileStore, s.config)
// userHandler.RegisterRoutes(api, "user-profile")
