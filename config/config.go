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

type Cors struct {
	AllowOrigins string `yaml:"allow_origin" env:"allow_origin" env-required:"true"`
}

type ClothesRoute struct {
	Product  string `yaml:"product" env:"product" env-required:"true"`
	Filter   string `yaml:"filter" env:"filter" env-required:"true"`
	Category string `yaml:"category" env:"category" env-required:"true"`
	Cart     string `yaml:"cart" env:"cart" env-required:"true"`
}

type ClothesStore struct {
	ProductList     string `yaml:"product_list" env:"product_list" env-required:"true"`
	ProductDetail   string `yaml:"product_detail" env:"product_detail" env-required:"true"`
	ProductBatch    string `yaml:"product_batch" env:"product_batch" env-required:"true"`
	ProductTemplate string `yaml:"product_template" env:"product_template" env-required:"true"`
	Category        string `yaml:"category" env:"category" env-required:"true"`
	Filter          string `yaml:"filter" env:"filter" env-required:"true"`
	FilterType      string `yaml:"fitler_type" env:"fitler_type" env-required:"true"`
	Cart            string `yaml:"cart" env:"cart" env-required:"true"`
}

type MasterRoute struct {
	Master string `yaml:"master" env:"master" env-required:"true"`
}

type MasterStore struct {
	Master string `yaml:"master" env:"master" env-required:"true"`
	State  string `yaml:"state" env:"state" env-required:"true"`
	City   string `yaml:"city" env:"city" env-required:"true"`
}

type Master struct {
	Routes       MasterRoute `yaml:"routes" env:"routes" env-required:"true"`
	DBCollection MasterStore `yaml:"db_collection" env:"db_collection" env-required:"true"`
}

type Clothes struct {
	Routes       ClothesRoute `yaml:"routes" env:"routes" env-required:"true"`
	DBCollection ClothesStore `yaml:"db_collection" env:"db_collection" env-required:"true"`
}

type AuthRoute struct {
	UserAccount string `yaml:"user_account" env:"user_account" env-required:"true"`
	UserProfile string `yaml:"user_profile" env:"user_profile" env-required:"true"`
	User        string `yaml:"user" env:"user" env-required:"true"`
	Owner       string `yaml:"owner" env:"owner" env-required:"true"`
	Bussiness   string `yaml:"bussiness" env:"bussiness" env-required:"true"`
}

type AuthStore struct {
	UserAccount string `yaml:"user_account" env:"user_account" env-required:"true"`
	UserProfile string `yaml:"user_profile" env:"user_profile" env-required:"true"`
	Owner       string `yaml:"owner" env:"owner" env-required:"true"`
	Bussiness   string `yaml:"bussiness" env:"bussiness" env-required:"true"`
}

type AuthClient struct {
	RedirectUrl string `yaml:"redirect_url" env:"redirect_url" env-required:"true"`
}

type Auth struct {
	Routes       AuthRoute  `yaml:"routes" env:"routes" env-required:"true"`
	DBCollection AuthStore  `yaml:"db_collection" env:"db_collection" env-required:"true"`
	Client       AuthClient `yaml:"client" env:"client" env-required:"true"`
}

type App struct {
	Clothes Clothes `yaml:"clothes" env:"clothes" env-required:"true"`
	Auth    Auth    `yaml:"auth" env:"auth" env-required:"true"`
	Master  Master  `yaml:"master" env:"master" env-required:"true"`
}

type Config struct {
	Env        string             `yaml:"env" env:"env" env-required:"true"`
	HTTPServer HTTPServer         `yaml:"http_server"`
	Storage    Storage            `yaml:"storage"`
	Secret     Secret             `yaml:"secrets"`
	AuthConfig AuthConfigProvider `yaml:"auth_config_provider"`
	App        App                `yaml:"app" env:"app" env-required:"true"`

	Cors Cors `yaml:"cors" env:"cors" env-required:"true"`

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
