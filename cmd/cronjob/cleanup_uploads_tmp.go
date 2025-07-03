package cronjob

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/robfig/cron/v3"
)

func StartCronJob() {
	c := cron.New()

	// Tambahkan cron yang jalan setiap jam 12 malam
	c.AddFunc("0 0 * * *", func() {
		fmt.Println("[CRON] Mulai bersihkan folder tmp:", time.Now().Format("2006-01-02 15:04:05"))
		cleanupTmpFolder()
	})

	c.Start()
}

func cleanupTmpFolder() {
	tmpDir := "/public/uploads/tmp"
	exclude := map[string]bool{
		"images": true,
		"docs":   true,
	}

	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		fmt.Println("Gagal membaca direktori:", err)
		return
	}

	for _, entry := range entries {
		name := entry.Name()
		if exclude[name] {
			continue
		}

		fullPath := filepath.Join(tmpDir, name)

		err := os.RemoveAll(fullPath)
		if err != nil {
			fmt.Printf("Gagal menghapus %s: %v\n", fullPath, err)
		} else {
			fmt.Printf("Berhasil menghapus: %s\n", fullPath)
		}
	}
}
