package config

import "os"

type TelegramConfig struct {
	BotToken      string
	BaseURL       string
	DefaultChatID string
}

func LoadTelegramConfig() TelegramConfig {
	return TelegramConfig{
		BotToken:      os.Getenv("TELEGRAM_BOT_TOKEN"),
		BaseURL:       "https://api.telegram.org",
		DefaultChatID: "888560906",
	}
}
