package responses

import "time"

type PasswordResetTokenResponse struct {
	ID *string `json:"id,omitempty"`
	Email *string `json:"email,omitempty"`
	Token *string `json:"token,omitempty"`
	ExpiredAt *time.Time `json:"expired_at,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
