package cronjob

import (
	"fmt"
	"os"
	"path/filepath"
)

func cleanupTmpFolder() {
	tmpDir := "./public/uploads/tmp"
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
		}
		fmt.Printf("Berhasil menghapus: %s\n", fullPath)
	}
	fmt.Println("Berhasil membersihkan folder tmp")
}
