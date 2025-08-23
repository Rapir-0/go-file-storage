package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "Filename required", http.StatusBadRequest)
		return
	}

	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	safePath := filepath.Join("../storage", filename)

	fmt.Printf("Downloading %s\n", safePath)

	absPath, _ := filepath.Abs(safePath)
	storageAbs, _ := filepath.Abs("../storage")

	if !strings.HasPrefix(absPath, storageAbs) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	http.ServeFile(w, r, safePath)
}
