package helpers

import (
	"path/filepath"
	"strconv"
	"time"
)

func GenerateNewFileName(oldName string) string {
	// Dapatkan tanggal dan waktu saat ini
	currentTime := time.Now()

	// Dapatkan timestamp dalam format yang unik
	timestamp := strconv.FormatInt(currentTime.UnixNano(), 10)

	// Dapatkan ekstensi file dari nama file lama
	ext := filepath.Ext(oldName)

	// Buat nama baru dengan format yang diinginkan dan tambahkan timestamp
	newName := "photo-" + timestamp + ext
	return newName
}
