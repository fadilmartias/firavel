package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        string    `gorm:"primarykey;not null"`
	Name      string    `gorm:"not null" faker:"name"`
	Email     string    `gorm:"unique;not null" faker:"email"`
	Password  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:NULL ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt time.Time `gorm:"type:timestamp;default:NULL"`
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
