package models

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID              string     `gorm:"primarykey;not null;size:7" json:"id,omitempty"`
	Name            string     `gorm:"not null;size:100" faker:"name" json:"name,omitempty"`
	Email           string     `gorm:"unique;not null;size:100" faker:"email" json:"email,omitempty"`
	Password        string     `gorm:"not null;size:100" json:"password,omitempty"`
	Role            string     `gorm:"not null;enum:admin,user;size:100" json:"role,omitempty"`
	EmailVerifiedAt *time.Time `gorm:"type:timestamp" json:"email_verified_at,omitempty"`
	CreatedAt       time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at,omitempty"`
	UpdatedAt       *time.Time `gorm:"type:timestamp;default:NULL ON UPDATE CURRENT_TIMESTAMP" json:"updated_at,omitempty"`
	DeletedAt       *time.Time `gorm:"type:timestamp;default:NULL" json:"deleted_at,omitempty"`
}

func generateShortID(length int) string {
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
		u.ID = generateShortID(7)
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
