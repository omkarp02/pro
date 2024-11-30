package config

import (
	"log"
	"os"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string `yaml:"address" env:"address" env-required:"true"`
}

type Storage struct {
	DBUrl  string `yaml:"db_url" env:"db_url" env-required:"true"`
	DBName string `yaml:"db_name" env:"db_name" env-required:"true"`
}

type Secret struct {
	AccessTokenPrivateKey  string `yaml:"access_token_private_key" env:"access_token_private_key" env-required:"true"`
	AccessTokenPublicKey   string `yaml:"access_token_public_key" env:"access_token_public_key" env-required:"true"`
	RefreshTokenPrivateKey string `yaml:"refresh_token_private_key" env:"refresh_token_private_key" env-required:"true"`
	RefreshTokenPublicKey  string `yaml:"refresh_token_public_key" env:"refresh_token_public_key" env-required:"true"`
	CookieEncryptionKey    string `yaml:"cookie_encryption_key" env:"cookie_encryption_key" env-required:"true"`
}

type Config struct {
	Env        string `yaml:"env" env:"env" env-required:"true"`
	HTTPServer `yaml:"http_server"`
	Storage    `yaml:"storage"`
	Secret     `yaml:"secrets"`
}

func MustLoad(configPath string) *Config {

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config file: %s", err.Error())
	}

	cfg.Secret.AccessTokenPrivateKey = strings.ReplaceAll(cfg.Secret.AccessTokenPrivateKey, "/\\n/g", "\n")
	cfg.Secret.AccessTokenPublicKey = strings.ReplaceAll(cfg.Secret.AccessTokenPublicKey, "/\\n/g", "\n")

	return &cfg

}
