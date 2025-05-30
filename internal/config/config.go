package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Port int `env:"PORT" env-default:"8080"`
}

func New() (*Config, error) {
	var config Config

	err := cleanenv.ReadEnv(&config)

	if err != nil {
		return &config, nil
	}
	return &config, nil
}
