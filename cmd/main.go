package main

import (
	"github.com/go-file-storage/handlers"
	"net/http"
)

func main() {
	//Uplaod handler
	http.HandleFunc("/api/upload", handlers.UploadHandler)

	//Download handler
	http.HandleFunc("/api/download", handlers.DownloadHandler)

	//Health check handler
	http.HandleFunc("/api/health", handlers.HealthHandler)

	http.ListenAndServe(":8080", nil)
}
