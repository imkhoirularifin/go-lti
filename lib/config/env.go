package config

import (
	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
)

type AppConfig struct {
	Port         string `env:"PORT" envDefault:"3000"`
	CanvasConfig CanvasConfig
	LtiConfig    CanvasLtiConfig
	ApiKeyConfig CanvasApiKeyConfig
	KeyConfig    KeyConfig
}

type CanvasConfig struct {
	Domain string `env:"CANVAS_DOMAIN"`
}

type CanvasLtiConfig struct {
	Issuer    string `env:"CANVAS_LTI_ISSUER"`
	JwkKid    string `env:"CANVAS_LTI_JWK_KID"`
	ClientId  string `env:"CANVAS_LTI_CLIENT_ID"`
	LaunchUrl string `env:"CANVAS_LTI_LAUNCH_URL"`
}

type CanvasApiKeyConfig struct {
	ClientId    string `env:"CANVAS_API_KEY_CLIENT_ID"`
	Secret      string `env:"CANVAS_API_KEY_SECRET"`
	RedirectUrl string `env:"CANVAS_API_KEY_REDIRECT_URL"`
}

type KeyConfig struct {
	PrivateKeyPath string `env:"PRIVATE_KEY_PATH"`
	PublicKeyPath  string `env:"PUBLIC_KEY_PATH"`
}

func Setup() (AppConfig, error) {
	var cfg AppConfig
	if err := env.Parse(&cfg); err != nil {
		return AppConfig{}, err
	}

	return cfg, nil
}
