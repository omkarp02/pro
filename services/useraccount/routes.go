package useraccount

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/services/middleware"
	services "github.com/omkarp02/pro/services/utils"
	"github.com/omkarp02/pro/types"
	"github.com/omkarp02/pro/utils"
)

type UserAccountStore interface {
	CreateUserAccount(CreateUserAccountType) (interface{}, error)
	GetUserAccountByEmail(string) (*UserAccount, error)
	GetUserAccount(query map[string]interface{}, project map[string]interface{}) (*UserAccount, error)
	GetUserFromRefreshToken(refreshToken string) (UserAccount, error)
	UpdateUserAccountById(id string, userAccount UserAccount) (bool, error)
	UpdateUserRefreshToken(userId string, action string, refreshToken string) error
	PullUserRefreshToken(refreshToken string) error
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

	var user CreateUserAccountType

	if err := services.ValidateBody(c, &user); err != nil {
		return err
	}

	hashedPassword, err := services.HashPassword(user.Password)
	if err != nil {
		utils.GenerateError(fiber.StatusInternalServerError, "error occured while generating password")
	}

	user.Password = hashedPassword

	id, err := h.store.CreateUserAccount(user)
	if err != nil {
		return err
	}

	slog.Debug("UserId", "id", fmt.Sprintf("%v", id))
	return utils.SendResponse(c, "User registered successfully", fiber.Map{"id": id}, 201)
}

func (h *Handler) login(c *fiber.Ctx) error {
	oldRefreshToken := c.Cookies(types.REFRESH_TOKEN_COOKIE)
	var userCred LoginUserAccountType

	if err := services.ValidateBody(c, &userCred); err != nil {
		return err
	}

	userAccount, err := h.store.GetUserAccountByEmail(userCred.Email)
	if err != nil {
		return err
	}

	if ok := services.CheckPasswordHash(userCred.Password, userAccount.Password); !ok {
		return utils.InvalidCredentails()
	}

	accessTokenPayload := types.ACCESS_TOKEN_PAYLOAD{
		ID: userAccount.ID.Hex(),
	}

	refreshTokenPayload := types.ACCESS_TOKEN_PAYLOAD{
		ID: userAccount.ID.Hex(),
	}

	accessToken, newRefreshToken, err := utils.GenerateRefreshAndAccessToken(accessTokenPayload, refreshTokenPayload, h.cfg)
	if err != nil {
		return err
	}

	var action string

	if len(oldRefreshToken) == 0 {
		action = "push"
	} else {
		if err := h.store.PullUserRefreshToken(oldRefreshToken); err != nil {
			if errors.Is(err, services.ErrDocumentNotFound) {
				action = "reinitialize"
			} else {
				return err
			}
		} else {
			action = "push"
		}

		clearCookie(c, types.REFRESH_TOKEN_COOKIE)
	}

	if err := h.store.UpdateUserRefreshToken(userAccount.ID.Hex(), action, newRefreshToken); err != nil {
		return err
	}

	updateRefreshTokenCookie(c, newRefreshToken)

	return utils.SendResponse(c, "User Logged In Succesfully", fiber.Map{"accessToken": accessToken}, 200)
}

func (h *Handler) handleRefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies(types.REFRESH_TOKEN_COOKIE)
	if len(refreshToken) == 0 {
		return utils.UnAuthorized("UnAuthorized")
	}

	clearCookie(c, types.REFRESH_TOKEN_COOKIE)

	_, err := h.store.GetUserFromRefreshToken(refreshToken)
	if errors.Is(err, services.ErrDocumentNotFound) {
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

	accessTokenPayload := types.ACCESS_TOKEN_PAYLOAD{
		ID: decodedUserData.ID,
	}

	refreshTokenPayload := types.REFRESH_TOKEN_PAYLOAD{
		ID: decodedUserData.ID,
	}

	accessToken, refreshToken, err := utils.GenerateRefreshAndAccessToken(accessTokenPayload, refreshTokenPayload, h.cfg)
	if err != nil {
		return err
	}

	if err := h.store.UpdateUserRefreshToken(decodedUserData.ID, "push", refreshToken); err != nil {
		return err
	}

	updateRefreshTokenCookie(c, refreshToken)

	return utils.SendResponse(c, "token generated successfully", fiber.Map{"accessToken": accessToken}, 200)

}

func (h *Handler) logout(c *fiber.Ctx) error {
	refreshToken := c.Cookies(types.REFRESH_TOKEN_COOKIE)
	if len(refreshToken) == 0 {
		return utils.UnAuthorized("UnAuthorized")
	}

	if err := h.store.PullUserRefreshToken(refreshToken); err != nil {
		clearCookie(c, types.REFRESH_TOKEN_COOKIE)
		return err
	}

	clearCookie(c, types.REFRESH_TOKEN_COOKIE)
	return utils.SendResponse(c, "User logged out successfully", fiber.Map{}, 200)
}

func updateRefreshTokenCookie(c *fiber.Ctx, refreshToken string) {
	c.Cookie(&fiber.Cookie{
		Name:    types.REFRESH_TOKEN_COOKIE,
		Value:   refreshToken,
		Expires: types.REFRESH_TOKEN_COOKIE_EXPIRY,
		// HTTPOnly: true,
		// Secure:   true,
		// SameSite: fiber.CookieSameSiteStrictMode,
	})
}

func clearCookie(c *fiber.Ctx, name string) {
	c.Cookie(&fiber.Cookie{
		Name:    name,
		Value:   "",
		Expires: time.Now().Add(-time.Hour),
	})
}
