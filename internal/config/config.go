package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Port int `env:"PORT" env-default:"8030"`

	ApiKey string `env:"APIKEY" env-default:""`

	ProxyUrl string `env:"PROXY_URL" env-default:""`
}

func New() (*Config, error) {
	var config Config

	err := cleanenv.ReadEnv(&config)

	if err != nil {
		return &config, nil
	}
	return &config, nil
}
