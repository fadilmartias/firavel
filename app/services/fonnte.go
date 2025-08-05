package services

import (
	"github.com/fadilmartias/firavel/config"
	"github.com/go-resty/resty/v2"
)

type FonnteSendMessageRequest struct {
	Target      *string `gorm:"type:text;not null" json:"target"`               // required
	Message     *string `gorm:"type:text" json:"message"`                       // optional
	URL         *string `gorm:"type:text" json:"url"`                           // optional
	Filename    *string `gorm:"type:varchar(255)" json:"filename"`              // optional
	Schedule    *int64  `gorm:"type:bigint" json:"schedule"`                    // optional
	Delay       *string `gorm:"type:varchar(50)" json:"delay"`                  // optional
	CountryCode *string `gorm:"type:varchar(10);default:62" json:"countryCode"` // optional
	Location    *string `gorm:"type:varchar(100)" json:"location"`              // optional
	Typing      *bool   `gorm:"default:false" json:"typing"`                    // optional
	Choices     *string `gorm:"type:text" json:"choices"`                       // optional
	Select      *string `gorm:"type:varchar(20)" json:"select"`                 // optional
	PollName    *string `gorm:"type:varchar(100)" json:"pollname"`              // optional
	File        *string `gorm:"type:text" json:"file"`                          // optional (path or identifier, not actual binary)
	ConnectOnly *bool   `gorm:"default:false" json:"connectOnly"`               // optional
	FollowUp    *int    `gorm:"type:int" json:"followup"`                       // optional
	Data        *string `gorm:"type:text" json:"data"`                          // optional
	Sequence    *bool   `gorm:"default:false" json:"sequence"`                  // optional
	Preview     *bool   `gorm:"default:true" json:"preview"`                    // optional
}

func FonnteSendMessage(body FonnteSendMessageRequest) (*resty.Response, error) {
	fonnteConfig := config.LoadFonnteConfig()
	client := resty.New().
		SetBaseURL(fonnteConfig.BaseURL).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fonnteConfig.Token)
	return client.R().
		SetBody(body).
		Post("/send")
}
