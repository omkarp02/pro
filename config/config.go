package config

import (
	"flag"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Google struct {
	ClientSecret string `yaml:"client_secret" env:"CLIENT_SECRET" env-required:"true"`
	ClientId     string `yaml:"client_id" env:"CLIENT_ID" env-required:"true"`
}

type HTTPServer struct {
	Addr string `yaml:"address" env:"address" env-required:"true"`
}

type Storage struct {
	DBUrl  string `yaml:"db_url" env:"DB_URL" env-required:"true"`
	DBName string `yaml:"db_name" env:"db_name" env-required:"true"`
}

type Secret struct {
	AccessTokenPrivateKey  string `yaml:"ACCESS_TOKEN_PRIVATE_KEY" env:"ACCESS_TOKEN_PRIVATE_KEY" env-required:"true"`
	AccessTokenPublicKey   string `yaml:"ACCESS_TOKEN_PUBLIC_KEY" env:"ACCESS_TOKEN_PUBLIC_KEY" env-required:"true"`
	RefreshTokenPrivateKey string `yaml:"REFRESH_TOKEN_PRIVATE_KEY" env:"REFRESH_TOKEN_PRIVATE_KEY" env-required:"true"`
	RefreshTokenPublicKey  string `yaml:"REFRESH_TOKEN_PUBLIC_KEY" env:"REFRESH_TOKEN_PUBLIC_KEY" env-required:"true"`
	CookieEncryptionKey    string `yaml:"COOKIE_ENCRYPTION_KEY" env:"COOKIE_ENCRYPTION_KEY" env-required:"true"`
	Google                 `yaml:"google"`
}

type AuthConfig struct {
	RedirectUrl  string `yaml:"redirect_url"`
	ProviderId   string `yaml:"provider_id" env-required:"true"`
	ProviderName string `yaml:"provider_name" env-required:"true"`
}

type AuthConfigProvider struct {
	Google AuthConfig `yaml:"google"`
	JWT    AuthConfig `yaml:"jwt"`
}

type Config struct {
	Env                    string `yaml:"env" env:"env" env-required:"true"`
	HTTPServer             `yaml:"http_server"`
	Storage                `yaml:"storage"`
	Secret                 `yaml:"secrets"`
	AuthConfig             AuthConfigProvider `yaml:"auth_config_provider"`
	getProviderIdByNameMap map[string]string
}

func (cfg *Config) GetProviderIdByName(name string) string {

	if value, exists := cfg.getProviderIdByNameMap[name]; exists {
		return value
	}
	panic("invalid provider name")
}

func MustLoad(configPath string) *Config {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if configPath == "" {
		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()

		configPath = *flags
		if configPath == "" {
			configPath = os.Getenv("CONFIG_PATH")
		}

		if configPath == "" {
			log.Fatal("Config path is not set")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config file: %s", err.Error())
	}

	setAdditionalConfig(&cfg)

	return &cfg

}

func setAdditionalConfig(cfg *Config) {
	cfg.Secret.AccessTokenPrivateKey = strings.ReplaceAll(cfg.Secret.AccessTokenPrivateKey, "/\\n/g", "\n")
	cfg.Secret.AccessTokenPublicKey = strings.ReplaceAll(cfg.Secret.AccessTokenPublicKey, "/\\n/g", "\n")

	v := reflect.ValueOf(cfg.AuthConfig)
	cfg.getProviderIdByNameMap = make(map[string]string, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i) // Value of the field
		// Ensure the field is of type AuthConfig
		if field.Kind() == reflect.Struct {
			// Access ProviderId and ProviderName
			providerId := field.FieldByName("ProviderId").String()
			providerName := field.FieldByName("ProviderName").String()

			cfg.getProviderIdByNameMap[providerName] = providerId
		}
	}

}
