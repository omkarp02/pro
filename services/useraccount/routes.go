package useraccount

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/services/middleware"
	services "github.com/omkarp02/pro/services/utils"
	"github.com/omkarp02/pro/types"
	"github.com/omkarp02/pro/utils"
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
	store UserAccountStore
	cfg   *config.Config
}

func NewHandler(store UserAccountStore, cfg *config.Config) *Handler {
	return &Handler{store: store, cfg: cfg}
}

func (h *Handler) RegisterRoutes(router fiber.Router, link string) {
	routeGrp := router.Group(link)

	routeGrp.Get("/handle-refresh-token", h.handleRefreshToken)
	routeGrp.Post("/register", h.registerUser)
	routeGrp.Post("/login", h.login)

	routeGrp.Use(middleware.VerifyToken(h.cfg))
	routeGrp.Get("/logout", h.logout)
}

func (h *Handler) registerUser(c *fiber.Ctx) error {

	var user CreateUserAccountBody

	if err := services.ValidateBody(c, &user); err != nil {
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
	if errors.Is(err, utils.ErrDocumentAlreadyExist) {
		return utils.StatusBadRequest("User already exist")
	} else if err != nil {
		return err
	}

	return utils.SendResponse(c, "User registered successfully", fiber.Map{"id": id}, 201)
}

func (h *Handler) login(c *fiber.Ctx) error {
	oldRefreshToken := c.Cookies(types.REFRESH_TOKEN_COOKIE)
	jwtProviderId := h.cfg.AuthConfig.JWT.ProviderId
	var userCred LoginUserAccountType

	if err := services.ValidateBody(c, &userCred); err != nil {
		return err
	}

	userAccount, err := h.store.GetUserAccountByEmail(userCred.Email)
	if errors.Is(err, utils.ErrDocumentNotFound) {
		return utils.StatusBadRequest("Invalid Credentials")
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
		return utils.InvalidCredentails()
	}

	if ok := services.CheckPasswordHash(userCred.Password, userAccount.PasswordHash); !ok {
		return utils.InvalidCredentails()
	}

	accessTokenPayload := services.CreateAccessTokenPayload(userId.Hex(), jwtProviderId)
	refreshTokenPayload := services.CreateRefreshTokenPayload(userId.Hex(), jwtProviderId)

	accessToken, newRefreshToken, err := utils.GenerateRefreshAndAccessToken(accessTokenPayload, refreshTokenPayload, h.cfg)
	if err != nil {
		return err
	}

	h.store.HandleRefreshTokenForLogin(userId.Hex(), newRefreshToken, oldRefreshToken)

	if len(oldRefreshToken) != 0 {
		services.ClearCookie(c, types.REFRESH_TOKEN_COOKIE)
	}

	services.UpdateRefreshTokenCookie(c, types.REFRESH_TOKEN_COOKIE, newRefreshToken)

	return utils.SendResponse(c, "User Logged In Succesfully", fiber.Map{"accessToken": accessToken}, 200)
}

func (h *Handler) handleRefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies(types.REFRESH_TOKEN_COOKIE)
	if len(refreshToken) == 0 {
		return utils.UnAuthorized("UnAuthorized")
	}

	services.ClearCookie(c, types.REFRESH_TOKEN_COOKIE)

	_, err := h.store.GetUserFromRefreshToken(refreshToken)
	if errors.Is(err, utils.ErrDocumentNotFound) {
		decodedUserData, err := utils.ValidateRefreshToken(refreshToken, h.cfg)
		if err != nil {
			return err
		}

		if err = h.store.UpdateUserRefreshToken(decodedUserData.ID, "empty", ""); err != nil {
			return err
		}

		return utils.UnAuthorized("Unauthorized")
	} else if err != nil {
		return err
	}

	decodedUserData, err := utils.ValidateRefreshToken(refreshToken, h.cfg)
	if err != nil {
		slog.Error("error while decoding token", "err", err.Error())
		h.store.UpdateUserRefreshToken(decodedUserData.ID, "empty", "")
		return utils.UnAuthorized("UnAuthorized")
	}

	accessTokenPayload := services.CreateAccessTokenPayload(decodedUserData.ID, decodedUserData.ProviderId)
	refreshTokenPayload := services.CreateRefreshTokenPayload(decodedUserData.ID, decodedUserData.ProviderId)

	accessToken, refreshToken, err := utils.GenerateRefreshAndAccessToken(accessTokenPayload, refreshTokenPayload, h.cfg)
	if err != nil {
		return err
	}

	if err := h.store.UpdateUserRefreshToken(decodedUserData.ID, "push", refreshToken); err != nil {
		return err
	}

	services.UpdateRefreshTokenCookie(c, types.REFRESH_TOKEN_COOKIE, refreshToken)

	return utils.SendResponse(c, "token generated successfully", fiber.Map{"accessToken": accessToken}, 200)

}

func (h *Handler) logout(c *fiber.Ctx) error {
	user := services.ValidateDataForAccessToken(c.Locals("user"))

	if user.ProviderId != h.cfg.AuthConfig.JWT.ProviderId {
		if err := goth_fiber.Logout(c); err != nil {
			return utils.InternalServerError()
		}
	}

	refreshToken := c.Cookies(types.REFRESH_TOKEN_COOKIE)
	if len(refreshToken) == 0 {
		return utils.UnAuthorized("UnAuthorized")
	}

	if err := h.store.PullUserRefreshToken(refreshToken); err != nil {
		services.ClearCookie(c, types.REFRESH_TOKEN_COOKIE)
		return err
	}

	services.ClearCookie(c, types.REFRESH_TOKEN_COOKIE)
	return utils.SendResponse(c, "User logged out successfully", fiber.Map{}, 200)
}
