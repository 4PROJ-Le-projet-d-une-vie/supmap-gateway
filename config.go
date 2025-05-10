package main

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Port           string `env:"SUPMAP_GATEWAY_PORT"`
	UsersHost      string `env:"SUPMAP_USERS_HOST"`
	UsersPort      string `env:"SUPMAP_USERS_PORT"`
	IncidentsHost  string `env:"SUPMAP_INCIDENTS_HOST"`
	IncidentsPort  string `env:"SUPMAP_INCIDENTS_PORT"`
	GisHost        string `env:"SUPMAP_GIS_HOST"`
	GisPort        string `env:"SUPMAP_GIS_PORT"`
	NavigationHost string `env:"SUPMAP_NAVIGATION_HOST"`
	NavigationPort string `env:"SUPMAP_NAVIGATION_PORT"`
}

func NewConfig() (*Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &cfg, nil
}
