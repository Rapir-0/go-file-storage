package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-file-storage/config"
	"github.com/go-file-storage/handlers"
)

func main() {
	// Парсим флаги командной строки
	var (
		configPath  = flag.String("config", "", "Path to configuration file (JSON or YAML)")
		genConfig   = flag.String("generate-config", "", "Generate example configuration file")
		showHelp    = flag.Bool("help", false, "Show help")
		showVersion = flag.Bool("version", false, "Show version")
	)
	flag.Parse()

	// Показываем версию
	if *showVersion {
		fmt.Println("Go File Storage v1.1.0")
		return
	}

	// Показываем помощь
	if *showHelp {
		flag.Usage()
		return
	}

	// Генерируем пример конфигурации
	if *genConfig != "" {
		cfg := &config.Config{
			Server: config.ServerConfig{
				Host: "0.0.0.0",
				Port: 8080,
				URL:  "http://localhost:8080",
			},
			Storage: config.StorageConfig{
				Path: "./storage",
			},
			Upload: config.UploadConfig{
				MaxFileSize: 10,
			},
		}

		if err := cfg.SaveExample(*genConfig); err != nil {
			log.Fatalf("Failed to generate config: %v", err)
		}

		fmt.Printf("Configuration file generated: %s\n", *genConfig)
		return
	}

	// Загружаем конфигурацию
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Создаем папку для хранения файлов
	if err := os.MkdirAll(cfg.Storage.Path, 0755); err != nil {
		log.Fatalf("Failed to create storage directory: %v", err)
	}

	// Инициализируем handlers с конфигурацией
	handlers.Init(cfg)

	// Настраиваем маршруты
	http.HandleFunc("/api/upload", handlers.UploadHandler)
	http.HandleFunc("/api/download", handlers.DownloadHandler)
	http.HandleFunc("/api/health", handlers.HealthHandler)

	// Запускаем сервер
	address := cfg.GetAddress()
	log.Printf("Starting server on %s", address)
	log.Printf("Storage path: %s", cfg.Storage.Path)
	log.Printf("Max file size: %d MB", cfg.Upload.MaxFileSize)
	log.Printf("Server URL: %s", cfg.Server.URL)

	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
