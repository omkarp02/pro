package useraccount

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

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
	CreateUser(ctx context.Context, payload CreateUserAccountModal) (string, error)
	GetUser(ctx context.Context, field string, value string) (UserAccount, error)
	UpdateUserRefreshToken(ctx context.Context, userId string, action string, refreshToken string) error
	PullUserRefreshToken(ctx context.Context, refreshToken string) error
	HandleRefreshTokenForLogin(ctx context.Context, userId string, refreshToken string, oldRefreshToken string) error
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
	routeGrp.Get("/user/logout", h.logout)

}

func (h *Handler) registerUser(c router.Context) error {

	ctx, cancel := createContext()
	defer cancel()

	var user TCreateUserAccount

	if err := h.validator.ValidateBody(c, &user); err != nil {
		return err
	}

	createUserAccountModal := CreateUserAccountModal{
		UserId:       user.UserId,
		PasswordHash: user.Password,
		Type:         constant.USERACCOUNT_TYPE_EMAIL,
		AuthProvider: []AuthProviderType{
			{
				Provider:   h.cfg.AuthConfig.JWT.ProviderName,
				ProviderID: h.cfg.AuthConfig.JWT.ProviderId,
			},
		},
		Role: []string{constant.ROLE_USER},
	}

	id, err := h.store.CreateUser(ctx, createUserAccountModal)
	if errors.Is(err, errutil.ErrDocumentAlreadyExist) {
		return errutil.AlreadyExist("User")
	} else if err != nil {
		return err
	}

	return utils.SendResponse(c, "User registered successfully", fiber.Map{"id": id}, 201)
}

func (h *Handler) login(c router.Context) error {
	ctx, cancel := createContext()
	defer cancel()

	oldRefreshToken := c.GetCookie(constant.REFRESH_TOKEN_COOKIE)
	jwtProviderId := h.cfg.AuthConfig.JWT.ProviderId
	var userCred LoginUserAccountType

	if err := h.validator.ValidateBody(c, &userCred); err != nil {
		return err
	}

	userAccount, err := h.store.GetUser(ctx, "userId", userCred.UserId)
	fmt.Println("1", err)
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
		return errutil.InternalServerError("not valid provider id")
	}

	if ok := helper.CheckPasswordHash(userCred.Password, userAccount.PasswordHash); !ok {
		return errutil.InvalidCredentails()
	}

	accessTokenPayload := helper.CreateAccessTokenPayload(userId.Hex(), jwtProviderId, userAccount.Role)
	refreshTokenPayload := helper.CreateRefreshTokenPayload(userId.Hex(), jwtProviderId, userAccount.Role)

	accessToken, newRefreshToken, err := utils.GenerateRefreshAndAccessToken(accessTokenPayload, refreshTokenPayload, h.cfg)
	if err != nil {
		return err
	}

	h.store.HandleRefreshTokenForLogin(ctx, userId.Hex(), newRefreshToken, oldRefreshToken)

	if len(oldRefreshToken) != 0 {
		helper.ClearCookie(c, constant.REFRESH_TOKEN_COOKIE)
	}

	helper.UpdateCookie(c, constant.REFRESH_TOKEN_COOKIE, newRefreshToken, constant.REFRESH_TOKEN_COOKIE_EXPIRY)

	fmt.Println("<<<<<<<<<<<<<<<<<<<< hre is the user logged in sucessuflly")

	return utils.SendResponse(c, "User Logged In Succesfully", fiber.Map{"accessToken": accessToken, "role": userAccount.Role, "userProfileId": userAccount.UserProfile}, 200)
}

func (h *Handler) handleRefreshToken(c router.Context) error {

	ctx, cancel := createContext()
	defer cancel()

	refreshToken := c.GetCookie(constant.REFRESH_TOKEN_COOKIE)

	if len(refreshToken) == 0 {
		return errutil.UnAuthorized("UnAuthorized")
	}

	helper.ClearCookie(c, constant.REFRESH_TOKEN_COOKIE)

	_, err := h.store.GetUser(ctx, "refresh_token", refreshToken)
	if errors.Is(err, errutil.ErrDocumentNotFound) {
		decodedUserData, err := utils.ValidateRefreshToken(refreshToken, h.cfg)
		if err != nil {
			return err
		}

		if err = h.store.UpdateUserRefreshToken(ctx, decodedUserData.ID, "empty", ""); err != nil {
			return err
		}

		return errutil.UnAuthorized("Unauthorized")
	} else if err != nil {
		return err
	}

	decodedUserData, err := utils.ValidateRefreshToken(refreshToken, h.cfg)
	if err != nil {
		slog.Error("error while decoding token", "err", err.Error())
		h.store.UpdateUserRefreshToken(ctx, decodedUserData.ID, "empty", "")
		return errutil.UnAuthorized("UnAuthorized")
	}

	accessTokenPayload := helper.CreateAccessTokenPayload(decodedUserData.ID, decodedUserData.ProviderId, decodedUserData.Role)
	refreshTokenPayload := helper.CreateRefreshTokenPayload(decodedUserData.ID, decodedUserData.ProviderId, decodedUserData.Role)

	accessToken, refreshToken, err := utils.GenerateRefreshAndAccessToken(accessTokenPayload, refreshTokenPayload, h.cfg)
	if err != nil {
		return err
	}

	if err := h.store.UpdateUserRefreshToken(ctx, decodedUserData.ID, "push", refreshToken); err != nil {
		return err
	}

	helper.UpdateCookie(c, constant.REFRESH_TOKEN_COOKIE, refreshToken, constant.REFRESH_TOKEN_COOKIE_EXPIRY)

	return utils.SendResponse(c, "token generated successfully", fiber.Map{"accessToken": accessToken}, 200)

}

func (h *Handler) authHandler(c router.Context) error {
	return goth_fiber.BeginAuthHandler(c.GetContext())
}

// func (h *Handler) redirectUrlHandler(c router.Context) error {

// 	ctx, cancel := createContext()
// 	defer cancel()

// 	provider := c.Params("provider")
// 	providerId := h.cfg.GetProviderIdByName(provider)

// 	oldRefreshToken := c.GetCookie(constant.REFRESH_TOKEN_COOKIE)
// 	user, err := goth_fiber.CompleteUserAuth(c.GetContext())

// 	fmt.Println(user.Name, user.FirstName, user.LastName)

// 	if err != nil {
// 		slog.Error("err while handling the redirect url", "err", err)
// 		return errutil.InternalServerError()
// 	}

// 	var id string
// 	userAccount, err := h.store.GetUser(ctx, "email", user.Email)
// 	id = userAccount.ID.Hex()

// 	if errors.Is(err, errutil.ErrDocumentNotFound) {
// 		createUserAccountModal := CreateUserAccountModal{
// 			Email: user.Email,
// 			AuthProvider: []AuthProviderType{
// 				{
// 					Provider:   provider,
// 					ProviderID: providerId,
// 				},
// 			},
// 		}

// 		profile := userprofile.CreateUserModel{
// 			FirstName: user.FirstName,
// 			LastName:  user.LastName,
// 			Email:     user.Email,
// 		}

// 		id, err = h.store.CreateUserProfileAndAccount(profile, createUserAccountModal)

// 		if err != nil {
// 			return err
// 		}
// 	} else if err != nil {
// 		return err
// 	}

// 	accessTokenPayload := helper.CreateAccessTokenPayload(id, providerId)
// 	refreshTokenPayload := helper.CreateRefreshTokenPayload(id, providerId)

// 	newAuthToken, newRefreshToken, err := utils.GenerateRefreshAndAccessToken(accessTokenPayload, refreshTokenPayload, h.cfg)
// 	if err != nil {
// 		return err
// 	}

// 	h.store.HandleRefreshTokenForLogin(ctx, id, newRefreshToken, oldRefreshToken)

// 	if len(oldRefreshToken) != 0 {
// 		helper.ClearCookie(c, constant.REFRESH_TOKEN_COOKIE)
// 	}

// 	helper.UpdateCookie(c, constant.REFRESH_TOKEN_COOKIE, newRefreshToken, constant.REFRESH_TOKEN_COOKIE_EXPIRY)

// 	return c.Redirect(h.cfg.App.Auth.Client.RedirectUrl+"?token="+newAuthToken, fiber.StatusFound)
// }

// func (h *Handler) create(c router.Context) error {
// 	decodedUserId := c.GetDecodedData().ID

// 	var user userprofile.TCreateUser

// 	// Parse the JSON body into the struct
// 	if err := h.validator.ValidateBody(c, &user); err != nil {
// 		return err
// 	}

// 	id, err := h.store.CreateUserProfile(user, decodedUserId)
// 	if err != nil {
// 		return err
// 	}

// 	return utils.SendResponse(c, "User created successfully", id, 201)
// }

func (h *Handler) logout(c router.Context) error {

	ctx, cancel := createContext()
	defer cancel()

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

	if err := h.store.PullUserRefreshToken(ctx, refreshToken); err != nil {
		helper.ClearCookie(c, constant.REFRESH_TOKEN_COOKIE)
		return err
	}

	helper.ClearCookie(c, constant.REFRESH_TOKEN_COOKIE)
	return utils.SendResponse(c, "User logged out successfully", fiber.Map{}, 200)
}

func createContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
