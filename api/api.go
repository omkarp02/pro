package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/services/auth"
	"github.com/omkarp02/pro/services/owner"
	"github.com/omkarp02/pro/services/useraccount"
	"github.com/omkarp02/pro/utils/errutil"
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

	fiberConfig := fiber.Config{
		ErrorHandler: errutil.ErrorHandler,
	}

	app := fiber.New(fiberConfig)

	app.Get("/asdf", func(c *fiber.Ctx) error {
		return errutil.ErrDocumentNotFound
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173/",
		AllowCredentials: true,
	}))

	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: s.config.Secret.CookieEncryptionKey,
	}))

	newRouter := router.NewFiberRouter(app)
	api := newRouter.Group("/api/v1")

	validator := validation.NewValidator()

	userAccountStore := useraccount.NewStore(s.db, "user_account")
	userAccountHandler := useraccount.NewHandler(userAccountStore, s.config, validator)
	userAccountHandler.RegisterRoutes(api, "user-account")

	authHandler := auth.NewHandler(s.config, userAccountStore)
	authHandler.RegisterRoutes(api, "auth")

	ownerRepo := owner.NewRepo(s.db, "owners")
	ownerService := owner.NewService(ownerRepo)
	ownerHandler := owner.NewHandler(ownerService, s.config, validator)
	ownerHandler.RegisterRoutes(api, "owner")

	// userProfileStore := userprofile.NewStore(s.db, "user_profile")
	// userHandler := userprofile.NewHandler(userProfileStore, s.config)
	// userHandler.RegisterRoutes(api, "user-profile")

	return app.Listen(s.addr)
}
