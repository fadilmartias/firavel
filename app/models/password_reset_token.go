package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type PasswordResetToken struct {
	Id        string    `gorm:"primarykey;not null"`
	Email     string    `gorm:"not null" faker:"email"`
	Token     string    `gorm:"not null" faker:"password"`
	ExpiredAt time.Time `gorm:"not null" faker:"date"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:NULL ON UPDATE CURRENT_TIMESTAMP"`
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
