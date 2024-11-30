package useraccount

import (
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/omkarp02/pro/config"
	services "github.com/omkarp02/pro/services/utils"
	"github.com/omkarp02/pro/types"
	"github.com/omkarp02/pro/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserAccountStore interface {
	CreateUserAccount(CreateUserAccountType) (interface{}, error)
	GetUserAccountByEmail(string) (*UserAccount, error)
	UpdateUserRefreshToken(userId bson.ObjectID, refreshToken string) error
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

	routeGrp.Post("/register", h.registerUser)
	routeGrp.Post("/login", h.login)
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

	accessToken, refreshToken, err := utils.GenerateRefreshAndAccessToken(userAccount, userAccount, h.cfg)
	if err != nil {
		return err
	}

	if err := h.store.UpdateUserRefreshToken(userAccount.ID, refreshToken); err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Expires:  types.REFRESH_TOKEN_EXPIRY,
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
	})

	return utils.SendResponse(c, "User Logged In Succesfully", fiber.Map{"accessToken": accessToken}, 200)
}
