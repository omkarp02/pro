package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/services/useraccount"
	"github.com/omkarp02/pro/services/userprofile"
	"github.com/omkarp02/pro/utils"
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

	fiberConfig := fiber.Config{
		ErrorHandler: utils.ErrorHandler,
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	app := fiber.New(fiberConfig)

	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: s.config.Secret.CookieEncryptionKey,
	}))

	api := app.Group("/api/v1")
	api.Use(injectValidator(validate))

	userProfileStore := userprofile.NewStore(s.db, "user_profile")
	userHandler := userprofile.NewHandler(userProfileStore, s.config)
	userHandler.RegisterRoutes(api, "user-profile")

	userAccountStore := useraccount.NewStore(s.db, "user_account")
	userAccountHandler := useraccount.NewHandler(userAccountStore, s.config)
	userAccountHandler.RegisterRoutes(api, "user-account")

	return app.Listen(s.addr)
}
