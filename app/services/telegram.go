package services

import (
	"github.com/fadilmartias/firavel/config"
	"github.com/go-resty/resty/v2"
)

func TelegramRequest(endpoint string, body map[string]any) (*resty.Response, error) {
	telegramConfig := config.LoadTelegramConfig()
	client := resty.New().
		SetBaseURL(telegramConfig.BaseURL+"/bot"+telegramConfig.BotToken).
		SetHeader("Content-Type", "application/json")
	return client.R().
		SetBody(body).
		Post(endpoint)
}

func TelegramSendMessage(text string) (*resty.Response, error) {
	body := map[string]any{
		"chat_id": config.LoadTelegramConfig().DefaultChatID,
		"text":    text,
	}
	return TelegramRequest("/sendMessage", body)
}
