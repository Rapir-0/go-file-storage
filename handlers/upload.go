package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if appConfig == nil {
		http.Error(w, "Server configuration not initialized", http.StatusInternalServerError)
		return
	}

	// Parse multipart form с динамическим лимитом
	maxSize := appConfig.GetMaxFileSizeBytes()
	err := r.ParseMultipartForm(maxSize)
	if err != nil {
		errMsg := fmt.Sprintf("File size exceeds %d MB limit", appConfig.Upload.MaxFileSize)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	// Get file from form
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Generate unique filename to avoid conflicts
	fileID := generateUniqueID()
	fileExt := filepath.Ext(handler.Filename)
	uniqueFilename := fileID + fileExt

	// Используем путь из конфигурации
	storagePath := appConfig.Storage.Path

	// Create destination file
	filePath := filepath.Join(storagePath, uniqueFilename)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy file content (streaming, better for large files)
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	// Prepare response
	response := map[string]interface{}{
		"success":       true,
		"file_id":       fileID,
		"original_name": handler.Filename,
		"download_url":  appConfig.Server.URL + "/api/download?filename=" + uniqueFilename,
		"size":          handler.Size,
	}

	// Set headers BEFORE WriteHeader
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode response
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		// Логируем ошибку, но статус уже отправлен
		println("Error encoding response:", err.Error())
	}
}

func generateUniqueID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
