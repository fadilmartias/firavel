package config

import "os"

type AppConfig struct {
	Name string
	Env  string
	Port string
}

func LoadAppConfig() AppConfig {
	return AppConfig{
		Name: os.Getenv("APP_NAME"),
		Env:  os.Getenv("APP_ENV"),
		Port: os.Getenv("APP_PORT"),
	}
}