package factories

import (
	"log"
	"math/rand"
	"time"

	"github.com/fadilmartias/firavel/app/models"
	"github.com/fadilmartias/firavel/app/utils"

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
	user.ID = utils.GenerateShortID(7)
	roles := []string{"admin", "user"}
	user.Role = roles[rand.Intn(len(roles))]
	user.Password = "password"
	createdAt := fakerDateBetween(time.Now().AddDate(-2, 0, 0), time.Now())
	user.CreatedAt = createdAt
	user.EmailVerifiedAt = nil
	user.UpdatedAt = nil
	user.DeletedAt = nil
	return user
}
