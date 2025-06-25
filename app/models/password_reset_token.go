package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type PasswordResetToken struct {
	ID        string     `gorm:"primarykey;size:7"`
	Email     string     `gorm:"not null;size:100" faker:"email"`
	Token     string     `gorm:"not null;size:255" faker:"password"`
	ExpiredAt *time.Time `gorm:"not null" faker:"date"`
	CreatedAt time.Time  `gorm:"not null"`
	UpdatedAt time.Time
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
