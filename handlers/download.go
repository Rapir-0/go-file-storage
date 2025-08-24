package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	if appConfig == nil {
		http.Error(w, "Server configuration not initialized", http.StatusInternalServerError)
		return
	}

	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "Filename required", http.StatusBadRequest)
		return
	}

	// Проверяем на path traversal атаки
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	// Используем путь из конфигурации
	storagePath := appConfig.Storage.Path
	safePath := filepath.Join(storagePath, filename)

	fmt.Printf("Downloading %s\n", safePath)

	// Проверяем что файл находится в папке storage
	absPath, _ := filepath.Abs(safePath)
	storageAbs, _ := filepath.Abs(storagePath)

	if !strings.HasPrefix(absPath, storageAbs) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Проверяем существование файла
	if _, err := os.Stat(safePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, safePath)
}
