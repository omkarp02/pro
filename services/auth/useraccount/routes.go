package useraccount

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/router"
	"github.com/omkarp02/pro/services/middleware"
	"github.com/omkarp02/pro/services/utils/helper"
	"github.com/omkarp02/pro/utils"
	"github.com/omkarp02/pro/utils/constant"
	"github.com/omkarp02/pro/utils/errutil"
	"github.com/omkarp02/pro/utils/validation"
	"github.com/shareed2k/goth_fiber"
)

type UserAccountStore interface {
	CreateUserAccount(CreateUserAccountModal) (string, error)
	GetUserAccountByEmail(string) (UserAccount, error)
	GetUserAccount(query map[string]interface{}, project map[string]interface{}) (*UserAccount, error)
	GetUserFromRefreshToken(refreshToken string) (UserAccount, error)
	UpdateUserAccountById(id string, userAccount UserAccount) (bool, error)
	UpdateUserRefreshToken(userId string, action string, refreshToken string) error
	PullUserRefreshToken(refreshToken string) error
	HandleRefreshTokenForLogin(userId string, refreshToken string, oldRefreshToken string) error
}

type Handler struct {
	store     UserAccountStore
	cfg       *config.Config
	validator *validation.Validator
}

func NewHandler(store UserAccountStore, cfg *config.Config, validator *validation.Validator) *Handler {
	return &Handler{store: store, cfg: cfg, validator: validator}
}

func (h *Handler) RegisterRoutes(router router.Router, link string) {
	routeGrp := router.Group(link)

	routeGrp.Get("/handle-refresh-token", h.handleRefreshToken)
	routeGrp.Post("/register", h.registerUser)
	routeGrp.Post("/login", h.login)

	routeGrp.Use(middleware.VerifyToken(h.cfg))
	routeGrp.Get("/logout", h.logout)
}

func (h *Handler) registerUser(c router.Context) error {

	var user CreateUserAccountBody

	if err := h.validator.ValidateBody(c, &user); err != nil {
		return err
	}

	createUserAccountModal := CreateUserAccountModal{
		Email:        user.Email,
		PasswordHash: user.Password,
		AuthProvider: []AuthProviderType{
			{
				Provider:   h.cfg.AuthConfig.JWT.ProviderName,
				ProviderID: h.cfg.AuthConfig.JWT.ProviderId,
			},
		},
	}

	id, err := h.store.CreateUserAccount(createUserAccountModal)
	if errors.Is(err, errutil.ErrDocumentAlreadyExist) {
		return errutil.AlreadyExist("User")
	} else if err != nil {
		return err
	}

	return utils.SendResponse(c, "User registered successfully", fiber.Map{"id": id}, 201)
}

func (h *Handler) login(c router.Context) error {
	oldRefreshToken := c.GetCookie(constant.REFRESH_TOKEN_COOKIE)
	jwtProviderId := h.cfg.AuthConfig.JWT.ProviderId
	var userCred LoginUserAccountType

	if err := h.validator.ValidateBody(c, &userCred); err != nil {
		return err
	}

	userAccount, err := h.store.GetUserAccountByEmail(userCred.Email)
	if errors.Is(err, errutil.ErrDocumentNotFound) {
		return errutil.StatusBadRequest("Invalid Credentials")
	} else if err != nil {
		return err
	}

	userId := userAccount.ID

	userHasJWTProvider := false
	for _, auth := range userAccount.AuthProvider {
		if auth.ProviderID == jwtProviderId {
			userHasJWTProvider = true
		}
	}

	if !userHasJWTProvider {
		return errutil.InvalidCredentails()
	}

	if ok := helper.CheckPasswordHash(userCred.Password, userAccount.PasswordHash); !ok {
		return errutil.InvalidCredentails()
	}

	accessTokenPayload := helper.CreateAccessTokenPayload(userId.Hex(), jwtProviderId)
	refreshTokenPayload := helper.CreateRefreshTokenPayload(userId.Hex(), jwtProviderId)

	accessToken, newRefreshToken, err := utils.GenerateRefreshAndAccessToken(accessTokenPayload, refreshTokenPayload, h.cfg)
	if err != nil {
		return err
	}

	h.store.HandleRefreshTokenForLogin(userId.Hex(), newRefreshToken, oldRefreshToken)

	if len(oldRefreshToken) != 0 {
		helper.ClearCookie(c, constant.REFRESH_TOKEN_COOKIE)
	}

	helper.UpdateCookie(c, constant.REFRESH_TOKEN_COOKIE, newRefreshToken, constant.REFRESH_TOKEN_COOKIE_EXPIRY)

	return utils.SendResponse(c, "User Logged In Succesfully", fiber.Map{"accessToken": accessToken}, 200)
}

func (h *Handler) handleRefreshToken(c router.Context) error {
	refreshToken := c.GetCookie(constant.REFRESH_TOKEN_COOKIE)
	if len(refreshToken) == 0 {
		return errutil.UnAuthorized("UnAuthorized")
	}

	helper.ClearCookie(c, constant.REFRESH_TOKEN_COOKIE)

	_, err := h.store.GetUserFromRefreshToken(refreshToken)
	if errors.Is(err, errutil.ErrDocumentNotFound) {
		decodedUserData, err := utils.ValidateRefreshToken(refreshToken, h.cfg)
		if err != nil {
			return err
		}

		if err = h.store.UpdateUserRefreshToken(decodedUserData.ID, "empty", ""); err != nil {
			return err
		}

		return errutil.UnAuthorized("Unauthorized")
	} else if err != nil {
		return err
	}

	decodedUserData, err := utils.ValidateRefreshToken(refreshToken, h.cfg)
	if err != nil {
		slog.Error("error while decoding token", "err", err.Error())
		h.store.UpdateUserRefreshToken(decodedUserData.ID, "empty", "")
		return errutil.UnAuthorized("UnAuthorized")
	}

	accessTokenPayload := helper.CreateAccessTokenPayload(decodedUserData.ID, decodedUserData.ProviderId)
	refreshTokenPayload := helper.CreateRefreshTokenPayload(decodedUserData.ID, decodedUserData.ProviderId)

	accessToken, refreshToken, err := utils.GenerateRefreshAndAccessToken(accessTokenPayload, refreshTokenPayload, h.cfg)
	if err != nil {
		return err
	}

	if err := h.store.UpdateUserRefreshToken(decodedUserData.ID, "push", refreshToken); err != nil {
		return err
	}

	helper.UpdateCookie(c, constant.REFRESH_TOKEN_COOKIE, refreshToken, constant.REFRESH_TOKEN_COOKIE_EXPIRY)

	return utils.SendResponse(c, "token generated successfully", fiber.Map{"accessToken": accessToken}, 200)

}

func (h *Handler) logout(c router.Context) error {
	user := c.GetDecodedData()

	if user.ProviderId != h.cfg.AuthConfig.JWT.ProviderId {
		if err := goth_fiber.Logout(c.GetContext()); err != nil {
			return errutil.InternalServerError()
		}
	}

	refreshToken := c.GetCookie(constant.REFRESH_TOKEN_COOKIE)
	if len(refreshToken) == 0 {
		return errutil.UnAuthorized("UnAuthorized")
	}

	if err := h.store.PullUserRefreshToken(refreshToken); err != nil {
		helper.ClearCookie(c, constant.REFRESH_TOKEN_COOKIE)
		return err
	}

	helper.ClearCookie(c, constant.REFRESH_TOKEN_COOKIE)
	return utils.SendResponse(c, "User logged out successfully", fiber.Map{}, 200)
}
