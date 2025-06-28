package models

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID              string  `gorm:"primarykey;size:7"`
	Name            string  `gorm:"not null;size:100" faker:"name"`
	Email           string  `gorm:"unique;not null;size:100" faker:"email"`
	Phone           string  `gorm:"unique;not null;size:15"`
	Password        string  `gorm:"not null;size:100" faker:"password"`
	Role            string  `gorm:"type:enum('admin','user');default:'user';not null"`
	RefreshToken    *string `gorm:"size:255"`
	EmailVerifiedAt *time.Time
	CreatedAt       time.Time `gorm:"not null"`
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

func GenerateID(length int) string {
	const shortIDChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = shortIDChars[seededRand.Intn(len(shortIDChars))]
	}
	return string(b)
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = GenerateID(7)
	}

	return
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
