package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/types"
	"github.com/omkarp02/pro/utils/constant"
)

type TokenGenerator interface {
	GenerateToken(payload interface{}, expiry time.Duration) (string, error)
	ValidateToken(token string) (interface{}, error)
}

type JWTTokenGenerator struct {
	publicKey  string
	privateKey string
}

func NewJWTTokenGenerator(publicKey string, privateKey string) *JWTTokenGenerator {

	return &JWTTokenGenerator{
		publicKey:  publicKey,
		privateKey: privateKey,
	}
}

func (gen *JWTTokenGenerator) GenerateToken(payload interface{}, expiry time.Duration) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(gen.privateKey))
	if err != nil {
		return "", err
	}

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["data"] = payload               // Our custom data.
	claims["exp"] = now.Add(expiry).Unix() // The expiration time after which the token must be disregarded.
	claims["iat"] = now.Unix()             // The time at which the token was issued.
	claims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (gen *JWTTokenGenerator) ValidateToken(token string) (interface{}, error) {
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(gen.publicKey))
	if err != nil {
		return nil, err
	}

	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return publicKey, nil
	})

	if err != nil {

		return nil, err
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, err
	}

	return claims["data"], nil
}

// TokenFactory creates a new TokenGenerator based on the token type
func TokenFactory(tokenType string, cfg *config.Config) (TokenGenerator, error) {
	secrets := cfg.Secret

	switch tokenType {
	case constant.ACCESS_TOKEN:
		return NewJWTTokenGenerator(secrets.AccessTokenPublicKey, secrets.AccessTokenPrivateKey), nil
	case constant.REFRESH_TOKEN:
		return NewJWTTokenGenerator(secrets.RefreshTokenPublicKey, secrets.RefreshTokenPrivateKey), nil
	default:
		return nil, fmt.Errorf("unsupported token type: %s", tokenType)
	}
}

func GenerateRefreshAndAccessToken(accessPayload interface{}, refreshPayload interface{}, cfg *config.Config) (string, string, error) {
	accessTokenGenerator, err := TokenFactory(constant.ACCESS_TOKEN, cfg)
	if err != nil {
		return "", "", err
	}

	accessToken, err := accessTokenGenerator.GenerateToken(accessPayload, time.Hour)
	if err != nil {
		return "", "", err
	}

	refreshTokenGenerator, err := TokenFactory(constant.REFRESH_TOKEN, cfg)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := refreshTokenGenerator.GenerateToken(refreshPayload, time.Hour)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func ValidateRefreshToken(refreshToken string, cfg *config.Config) (types.REFRESH_TOKEN_PAYLOAD, error) {
	refreshTokenPayload := types.REFRESH_TOKEN_PAYLOAD{}

	refreshTokenGenerator, err := TokenFactory(constant.REFRESH_TOKEN, cfg)
	if err != nil {
		return refreshTokenPayload, err
	}

	claimsMap, err := refreshTokenGenerator.ValidateToken(refreshToken)
	if err != nil {
		return refreshTokenPayload, err
	}

	decodedUserData, err := GetUserDataFromRefreshClaimsData(claimsMap)
	if err != nil {
		return refreshTokenPayload, err
	}

	return decodedUserData, nil
}
