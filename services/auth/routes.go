package auth

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/services/useraccount"
	"github.com/omkarp02/pro/services/utils/helper"
	"github.com/omkarp02/pro/utils"
	"github.com/omkarp02/pro/utils/constant"
	"github.com/omkarp02/pro/utils/errutil"
	"github.com/shareed2k/goth_fiber"
)

type UserAccountStore interface {
	CreateUserAccount(user useraccount.CreateUserAccountModal) (string, error)
	HandleRefreshTokenForLogin(userId string, refreshToken string, oldRefreshToken string) error
	GetUserAccountByEmail(email string) (useraccount.UserAccount, error)
}

type Handler struct {
	cfg *config.Config
	UserAccountStore
}

func NewHandler(cfg *config.Config, userAccountStore UserAccountStore) *Handler {
	return &Handler{cfg: cfg, UserAccountStore: userAccountStore}
}

func (h *Handler) RegisterRoutes(router router.Router, link string) {
	h.RegisterProviders()
	routeGrp := router.Group(link)

	routeGrp.Get("/:provider", h.authHandler)
	router.Get("/auth/:provider/callback", h.redirectUrlHandler)
}

func (h *Handler) RegisterProviders() {
	googleAuthConfig := h.cfg.AuthConfig.Google
	googleSecret := h.cfg.Secret.Google

	goth.UseProviders(
		google.New(googleSecret.ClientId, googleSecret.ClientSecret, googleAuthConfig.RedirectUrl),
	)
}

func (h *Handler) authHandler(c router.Context) error {
	return goth_fiber.BeginAuthHandler(c.GetContext())
}

func (h *Handler) redirectUrlHandler(c router.Context) error {

	provider := c.Params("provider")
	providerId := h.cfg.GetProviderIdByName(provider)

	oldRefreshToken := c.GetCookie(constant.REFRESH_TOKEN_COOKIE)
	user, err := goth_fiber.CompleteUserAuth(c.GetContext())
	if err != nil {
		slog.Error("err while handling the redirect url", "err", err)
		return errutil.InternalServerError()
	}

	slog.Info("Info about user", "user", user)

	var id string
	userAccount, err := h.GetUserAccountByEmail(user.Email)
	id = userAccount.ID.Hex()
	if errors.Is(err, errutil.ErrDocumentNotFound) {
		createUserAccountModal := useraccount.CreateUserAccountModal{
			Email: user.Email,
			AuthProvider: []useraccount.AuthProviderType{
				{
					Provider:   provider,
					ProviderID: providerId,
				},
			},
		}

		id, err = h.CreateUserAccount(createUserAccountModal)
	} else if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	accessTokenPayload := helper.CreateAccessTokenPayload(id, providerId)
	refreshTokenPayload := helper.CreateRefreshTokenPayload(id, providerId)

	_, newRefreshToken, err := utils.GenerateRefreshAndAccessToken(accessTokenPayload, refreshTokenPayload, h.cfg)
	if err != nil {
		return err
	}

	h.HandleRefreshTokenForLogin(id, newRefreshToken, oldRefreshToken)

	if len(oldRefreshToken) != 0 {
		helper.ClearCookie(c, constant.REFRESH_TOKEN_COOKIE)
	}

	helper.UpdateCookie(c, constant.REFRESH_TOKEN_COOKIE, newRefreshToken, constant.REFRESH_TOKEN_COOKIE_EXPIRY)

	return c.Redirect("http://localhost:5173/success", fiber.StatusFound)
}
