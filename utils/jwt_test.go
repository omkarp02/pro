package utils

import (
	"log/slog"
	"testing"
	"time"

	"github.com/omkarp02/pro/config"
	"github.com/omkarp02/pro/types"
)

func generateConfig() *config.Config {
	cfg := config.MustLoad("../config/local.yaml")
	return cfg
}

func TestJwt(t *testing.T) {
	cfg := generateConfig()
	tokenGenerator, err := TokenFactory(types.ACCESS_TOKEN, cfg)
	if err != nil {
		t.Error(err)
	}
	token, err := tokenGenerator.GenerateToken("alsdkfjasd", time.Hour)
	if err != nil {
		t.Error(err)
	}

	data, err := tokenGenerator.ValidateToken(token)
	if err != nil {
		t.Error(err)
	}

	slog.Debug("data", "data", data)
}
