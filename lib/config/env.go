package config

import (
	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
)

type AppConfig struct {
	Port      string `env:"PORT" envDefault:"3000"`
	LtiConfig LtiConfig
	KeyConfig KeyConfig
}

type LtiConfig struct {
	LaunchUrl    string `env:"LTI_LAUNCH_URL"`
	ClientId     string `env:"LTI_CLIENT_ID"`
	CanvasDomain string `env:"CANVAS_DOMAIN"`
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
