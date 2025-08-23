package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const ServerIP = "http://localhost:8080" // Добавил http://

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form (10MB limit)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "File size exceeds 10MB limit", http.StatusBadRequest)
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

	// Create storage directory if it doesn't exist
	os.MkdirAll("../storage", 0755)

	// Create destination file
	dst, err := os.Create("../storage/" + uniqueFilename)
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
		"download_url":  ServerIP + "/api/download?filename=" + uniqueFilename,
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
