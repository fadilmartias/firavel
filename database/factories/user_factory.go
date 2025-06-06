package factories

import (
	"goravel/app/models"
	"log"
	"math/rand"
	"time"

	"github.com/bxcodec/faker/v3"
)

func fakerDateBetween(start, end time.Time) time.Time {
	diff := end.Sub(start)
	// Ambil waktu acak antara start dan end
	randSeconds := time.Duration(rand.Int63n(int64(diff)))
	return start.Add(randSeconds)
}

// NewUser membuat instance User baru dengan data palsu tanpa menyimpannya.
func NewUser() models.User {
	var user models.User
	err := faker.FakeData(&user)
	if err != nil {
		log.Printf("Error faking user data: %v", err)
	}
	// Password tidak di-hash di sini, seeder yang akan melakukannya
	user.Password = "password"
	user.CreatedAt = fakerDateBetween(time.Now().AddDate(-2, 0, 0), time.Now())
	user.UpdatedAt = fakerDateBetween(user.CreatedAt, time.Now())
	user.DeletedAt = fakerDateBetween(user.CreatedAt, time.Now().AddDate(1, 0, 0))
	return user
}
