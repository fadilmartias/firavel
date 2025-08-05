package config

import "os"

type FonnteConfig struct {
	BaseURL string
	Token   string
}

func LoadFonnteConfig() FonnteConfig {
	return FonnteConfig{
		BaseURL: "https://api.fonnte.com",
		Token:   os.Getenv("FONNTE_TOKEN"),
	}
}
