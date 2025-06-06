package factories

import (
	"goravel/app/models"
	"log"

	"github.com/bxcodec/faker/v3"
)

// NewUser membuat instance User baru dengan data palsu tanpa menyimpannya.
func NewUser() models.User {
	var user models.User
	err := faker.FakeData(&user)
	if err != nil {
		log.Printf("Error faking user data: %v", err)
	}
	// Password tidak di-hash di sini, seeder yang akan melakukannya
	user.Password = "password"
	return user
}
