package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id              string    `gorm:"primarykey;not null;size:7" json:"id"`
	Name            string    `gorm:"not null;size:100" faker:"name" json:"name"`
	Email           string    `gorm:"unique;not null;size:100" faker:"email" json:"email"`
	Password        string    `gorm:"not null;size:100" json:"password"`
	Role            string    `gorm:"not null;size:100" json:"role"`
	EmailVerifiedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"email_verified_at"`
	CreatedAt       time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time `gorm:"type:timestamp;default:NULL ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt       time.Time `gorm:"type:timestamp;default:NULL" json:"deleted_at"`
}

// HashPassword mengenkripsi password sebelum disimpan
func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

// CheckPassword memverifikasi password
func (u *User) CheckPassword(providedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(providedPassword))
}
