package cronjob

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

func StartCronJob(db *gorm.DB) {
	c := cron.New()

	// Tambahkan cron yang jalan setiap jam 12 malam
	c.AddFunc("0 0 * * *", func() {
		fmt.Println("[CRON] Mulai bersihkan folder tmp:", time.Now().Format("2006-01-02 15:04:05"))
		cleanupTmpFolder()
	})

	c.Start()
}
