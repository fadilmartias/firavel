package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type PasswordResetToken struct {
	ID        string    `gorm:"primarykey;size:7"`
	Email     string    `gorm:"not null;size:100" faker:"email"`
	Token     string    `gorm:"not null;size:255" faker:"password"`
	ExpiredAt time.Time `gorm:"not null" faker:"date"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time
}

func (p *PasswordResetToken) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID = GenerateID(7)
	}

	return
}

// HashToken mengenkripsi token sebelum disimpan
func (p *PasswordResetToken) HashToken(token *string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*token), 14)
	if err != nil {
		return err
	}
	p.Token = string(bytes)
	return nil
}

// CheckToken memverifikasi token
func (p *PasswordResetToken) CheckToken(providedToken string) error {
	return bcrypt.CompareHashAndPassword([]byte(p.Token), []byte(providedToken))
}
