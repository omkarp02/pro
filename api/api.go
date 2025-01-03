package api

import (
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/services/clothes/filter"
	"github.com/omkarp02/pro/services/clothes/product"
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

	return api.Listen(s.addr)
}

func setUpClothesApp(curDb *db.Database, cfg *config.Config, validator *validation.Validator, api router.Router) {

	filter.Intialize(curDb, cfg, validator, api)
	product.Intialize(curDb, cfg, validator, api)
}

// userAccountStore := useraccount.NewStore(s.db, "user_account")
// userAccountHandler := useraccount.NewHandler(userAccountStore, s.config, validator)
// userAccountHandler.RegisterRoutes(api, "user-account")

// authHandler := auth.NewHandler(s.config, userAccountStore)
// authHandler.RegisterRoutes(api, "auth")

// ownerRepo := owner.NewRepo(s.db, "owners")
// ownerService := owner.NewService(ownerRepo)
// ownerHandler := owner.NewHandler(ownerService, s.config, validator)
// ownerHandler.RegisterRoutes(api, "owner")

// userProfileStore := userprofile.NewStore(s.db, "user_profile")
// userHandler := userprofile.NewHandler(userProfileStore, s.config)
// userHandler.RegisterRoutes(api, "user-profile")
