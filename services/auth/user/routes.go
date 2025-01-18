package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/services/auth/useraccount"
	"github.com/omkarp02/pro/services/auth/userprofile"
	"github.com/omkarp02/pro/services/middleware"
	"github.com/omkarp02/pro/services/utils/helper"
	"github.com/omkarp02/pro/utils"
	"github.com/omkarp02/pro/utils/constant"
	"github.com/omkarp02/pro/utils/errutil"
	"github.com/omkarp02/pro/utils/validation"
	"github.com/shareed2k/goth_fiber"
)

type UserService interface {
	CreateUserProfileAndAccount(ctx context.Context, userprofile userprofile.CreateUserModel, useraccount useraccount.CreateUserAccountModal) (string, error)
	CreateUserProfile(ctx context.Context, paylaod userprofile.TCreateUser, useraccountId string) (string, error)
	CreateOwnerAndAccount(ctx context.Context, ownerPayload CreateOwnerAndAccountModel, userId string) (string, error)
	GetUser(ctx context.Context, field string, value string) (useraccount.UserAccount, error)
	HandleRefreshTokenForLogin(ctx context.Context, userId string, refreshToken string, oldRefreshToken string) error
}

type Handler struct {
	service   UserService
	cfg       *config.Config
	validator *validation.Validator
}

func NewHandler(service UserService, cfg *config.Config, validator *validation.Validator) *Handler {
	return &Handler{service: service, cfg: cfg, validator: validator}
}

func (h *Handler) RegisterRoutes(router router.Router, link string) {
	h.RegisterProviders()
	routeGrp := router.Group(link)

	routeGrp.Get("/sso/:provider", h.authHandler)
	routeGrp.Get("/sso/:provider/callback", h.redirectUrlHandler)

	routeGrp.Use(middleware.VerifyToken(h.cfg))

	routeGrp.Post("/create-user-profile", h.createUserProfile)
	routeGrp.Post("/create-owner", h.createOwner)
}

func (h *Handler) RegisterProviders() {
	googleAuthConfig := h.cfg.AuthConfig.Google
	googleSecret := h.cfg.Secret.Google

	goth.UseProviders(
		google.New(googleSecret.ClientId, googleSecret.ClientSecret, googleAuthConfig.RedirectUrl, "profile", "email"),
	)
}

func (h *Handler) createUserProfile(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	decodedUserId := c.GetDecodedData().ID

	var user userprofile.TCreateUser

	// Parse the JSON body into the struct
	if err := h.validator.ValidateBody(c, &user); err != nil {
		return err
	}

	id, err := h.service.CreateUserProfile(ctx, user, decodedUserId)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "User created successfully", id, 201)
}

func (h *Handler) authHandler(c router.Context) error {
	return goth_fiber.BeginAuthHandler(c.GetContext())
}

func (h *Handler) redirectUrlHandler(c router.Context) error {

	ctx, cancel := createContext()
	defer cancel()

	role := []string{constant.ROLE_ADMIN}
	provider := c.Params("provider")
	providerId := h.cfg.GetProviderIdByName(provider)

	oldRefreshToken := c.GetCookie(constant.REFRESH_TOKEN_COOKIE)
	user, err := goth_fiber.CompleteUserAuth(c.GetContext())

	fmt.Println(user.Name, user.FirstName, user.LastName)

	if err != nil {
		slog.Error("err while handling the redirect url", "err", err)
		return errutil.InternalServerError()
	}

	var id string
	userAccount, err := h.service.GetUser(ctx, "email", user.Email)
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
			Role: role,
		}

		profile := userprofile.CreateUserModel{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		}

		id, err = h.service.CreateUserProfileAndAccount(ctx, profile, createUserAccountModal)

		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	accessTokenPayload := helper.CreateAccessTokenPayload(id, providerId, role)
	refreshTokenPayload := helper.CreateRefreshTokenPayload(id, providerId, role)

	newAuthToken, newRefreshToken, err := utils.GenerateRefreshAndAccessToken(accessTokenPayload, refreshTokenPayload, h.cfg)
	if err != nil {
		return err
	}

	h.service.HandleRefreshTokenForLogin(ctx, id, newRefreshToken, oldRefreshToken)

	if len(oldRefreshToken) != 0 {
		helper.ClearCookie(c, constant.REFRESH_TOKEN_COOKIE)
	}

	helper.UpdateCookie(c, constant.REFRESH_TOKEN_COOKIE, newRefreshToken, constant.REFRESH_TOKEN_COOKIE_EXPIRY)

	return c.Redirect(h.cfg.App.Auth.Client.RedirectUrl+"?token="+newAuthToken, fiber.StatusFound)
}

func (h *Handler) createOwner(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	decodedUserData := c.GetDecodedData()

	var owner TCreateOwner

	// Parse the JSON body into the struct
	if err := h.validator.ValidateBody(c, &owner); err != nil {
		return err
	}

	createOwnerData := CreateOwnerAndAccountModel{
		Name:         owner.Name,
		Email:        owner.Email,
		FirstName:    owner.FirstName,
		Password:     owner.Password,
		LastName:     owner.LastName,
		ProviderId:   h.cfg.AuthConfig.JWT.ProviderId,
		ProviderName: h.cfg.AuthConfig.Google.ProviderName,
		DateOfBirth:  owner.DateOfBirth,
		MobileNo:     owner.MobileNo,
		Gender:       owner.Gender,
	}

	id, err := h.service.CreateOwnerAndAccount(ctx, createOwnerData, decodedUserData.ID)
	if err != nil {
		return err
	}

	return utils.SendResponse(c, "User created successfully", id, 201)
}

func createContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
