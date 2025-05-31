package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Port int `env:"PORT" env-default:"8030"`

	ApiKey string `env:"APIKEY" env-default:"AIzaSyCRx6WlYdMfrOx8VYbYtrpfKPHYFZsgHbY"`

	ProxyUrl string `env:"PROXY_URL" env-default:"http://WCBDDLZI:HZ50YX0E@185.213.249.50:45596"`
}

func New() (*Config, error) {
	var config Config

	err := cleanenv.ReadEnv(&config)

	if err != nil {
		return &config, nil
	}
	return &config, nil
}
