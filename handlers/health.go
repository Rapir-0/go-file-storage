package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
		return
	}

	if appConfig == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Configuration not initialized"))
		return
	}

	// Проверяем доступность папки storage
	storageInfo := "OK"
	if _, err := os.Stat(appConfig.Storage.Path); os.IsNotExist(err) {
		storageInfo = "Storage directory not found"
	}

	// Подготавливаем ответ
	health := map[string]interface{}{
		"status":           "healthy",
		"version":          "1.1.0",
		"storage_path":     appConfig.Storage.Path,
		"storage_status":   storageInfo,
		"max_file_size_mb": appConfig.Upload.MaxFileSize,
		"server_url":       appConfig.Server.URL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(health); err != nil {
		fmt.Printf("Error encoding health response: %v\n", err)
	}
}
