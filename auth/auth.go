package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"github.com/omkarp02/pro/config"
	"github.com/shareed2k/goth_fiber"
)

func NewAuth(cfg *config.Config, router fiber.Router) {
	googleClientId := cfg.Secret.Google.ClientId
	googleClientSecret := cfg.Secret.Google.ClientSecret

	goth.UseProviders(
		google.New(googleClientId, googleClientSecret, "http://localhost:8080/auth/callback/google"),
	)

	router.Get("/login/:provider", goth_fiber.BeginAuthHandler)
	router.Get("/auth/callback/:provider", func(ctx *fiber.Ctx) error {
		user, err := goth_fiber.CompleteUserAuth(ctx)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(user.RefreshToken, "<<<<<<<< here is refresh token")

		data, err := refreshAccessToken(googleClientId, googleClientSecret, user.RefreshToken)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(data.AccessToken, data.ExpiresIn)

		return ctx.SendString(user.Email)
	})

	router.Get("/logout", func(ctx *fiber.Ctx) error {
		if err := goth_fiber.Logout(ctx); err != nil {
			log.Fatal(err)
		}

		return ctx.SendString("logout")
	})

}

const googleTokenEndpoint = "https://oauth2.googleapis.com/token"

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	IDToken      string `json:"id_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"` // Usually not returned again
}

func refreshAccessToken(clientID, clientSecret, refreshToken string) (*RefreshTokenResponse, error) {
	// Prepare the request payload
	payload := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"refresh_token": refreshToken,
		"grant_type":    "refresh_token",
	}

	fmt.Println(refreshToken)

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Make the HTTP POST request
	resp, err := http.Post(googleTokenEndpoint, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}

	defer resp.Body.Close()

	fmt.Println(resp, "<<<<<<<<<<<<<<< rposne")

	// Read and parse the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Println(string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}

	var tokenResponse RefreshTokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &tokenResponse, nil
}
