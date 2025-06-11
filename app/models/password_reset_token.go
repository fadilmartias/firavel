package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type PasswordResetToken struct {
	Id        string     `gorm:"primarykey;not null;size:7" json:"id,omitempty"`
	Email     string     `gorm:"not null;size:100" faker:"email" json:"email,omitempty"`
	Token     string     `gorm:"not null;size:255" faker:"password" json:"token,omitempty"`
	ExpiredAt *time.Time `gorm:"not null" faker:"date" json:"expired_at,omitempty"`
	CreatedAt *time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at,omitempty"`
	UpdatedAt *time.Time `gorm:"type:timestamp;default:NULL ON UPDATE CURRENT_TIMESTAMP" json:"updated_at,omitempty"`
}

// HashToken mengenkripsi token sebelum disimpan
func (u *PasswordResetToken) HashToken(token string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(token), 14)
	if err != nil {
		return err
	}
	u.Token = string(bytes)
	return nil
}

// CheckToken memverifikasi token
func (u *PasswordResetToken) CheckToken(providedToken string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Token), []byte(providedToken))
}
