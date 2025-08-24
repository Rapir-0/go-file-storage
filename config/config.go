package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config структура для всех настроек приложения
type Config struct {
	Server  ServerConfig  `json:"server" yaml:"server"`
	Storage StorageConfig `json:"storage" yaml:"storage"`
	Upload  UploadConfig  `json:"upload" yaml:"upload"`
}

// ServerConfig настройки сервера
type ServerConfig struct {
	Host string `json:"host" yaml:"host"`
	Port int    `json:"port" yaml:"port"`
	URL  string `json:"url" yaml:"url"` // Полный URL для формирования ссылок
}

// StorageConfig настройки хранилища
type StorageConfig struct {
	Path string `json:"path" yaml:"path"`
}

// UploadConfig настройки загрузки файлов
type UploadConfig struct {
	MaxFileSize int64 `json:"max_file_size_mb" yaml:"max_file_size_mb"` // в мегабайтах
}

// Load загружает конфигурацию из файла или переменных окружения
func Load(configPath string) (*Config, error) {
	cfg := &Config{
		// Значения по умолчанию
		Server: ServerConfig{
			Host: "localhost",
			Port: 8080,
			URL:  "http://localhost:8080",
		},
		Storage: StorageConfig{
			Path: "./storage",
		},
		Upload: UploadConfig{
			MaxFileSize: 10, // 10 MB
		},
	}

	// Пытаемся загрузить из файла
	if configPath != "" {
		if err := cfg.loadFromFile(configPath); err != nil {
			return nil, fmt.Errorf("failed to load config from file: %w", err)
		}
	}

	// Переопределяем из переменных окружения
	cfg.loadFromEnv()

	// Валидируем конфигурацию
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// loadFromFile загружает конфигурацию из JSON или YAML файла
func (c *Config) loadFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Файл не существует - это нормально
		}
		return err
	}
	defer file.Close()

	// Определяем тип файла по расширению
	ext := filepath.Ext(path)
	switch strings.ToLower(ext) {
	case ".yaml", ".yml":
		decoder := yaml.NewDecoder(file)
		return decoder.Decode(c)
	case ".json":
		decoder := json.NewDecoder(file)
		return decoder.Decode(c)
	default:
		// По умолчанию пробуем JSON
		decoder := json.NewDecoder(file)
		return decoder.Decode(c)
	}
}

// loadFromEnv загружает настройки из переменных окружения
func (c *Config) loadFromEnv() {
	// Server config
	if host := os.Getenv("SERVER_HOST"); host != "" {
		c.Server.Host = host
	}

	if port := os.Getenv("SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			c.Server.Port = p
		}
	}

	if url := os.Getenv("SERVER_URL"); url != "" {
		c.Server.URL = url
	} else {
		// Автоматически формируем URL
		c.Server.URL = fmt.Sprintf("http://%s:%d", c.Server.Host, c.Server.Port)
	}

	// Storage config
	if path := os.Getenv("STORAGE_PATH"); path != "" {
		c.Storage.Path = path
	}

	// Upload config
	if maxSize := os.Getenv("MAX_FILE_SIZE"); maxSize != "" {
		// Поддерживаем форматы: "10MB", "50mb", "1GB", "500"
		size, err := parseFileSize(maxSize)
		if err == nil {
			c.Upload.MaxFileSize = size
		}
	}
}

// parseFileSize парсит размер файла из строки (10MB, 1GB, etc.)
func parseFileSize(size string) (int64, error) {
	size = strings.ToUpper(strings.TrimSpace(size))

	// Если только число - считаем что мегабайты
	if num, err := strconv.ParseInt(size, 10, 64); err == nil {
		return num, nil
	}

	// Парсим с суффиксом
	var multiplier int64 = 1
	if strings.HasSuffix(size, "KB") {
		multiplier = 1
		size = strings.TrimSuffix(size, "KB")
	} else if strings.HasSuffix(size, "MB") {
		multiplier = 1
		size = strings.TrimSuffix(size, "MB")
	} else if strings.HasSuffix(size, "GB") {
		multiplier = 1024
		size = strings.TrimSuffix(size, "GB")
	}

	num, err := strconv.ParseInt(size, 10, 64)
	if err != nil {
		return 0, err
	}

	return num * multiplier, nil
}

// validate проверяет правильность конфигурации
func (c *Config) validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if c.Storage.Path == "" {
		return fmt.Errorf("storage path cannot be empty")
	}

	if c.Upload.MaxFileSize <= 0 {
		return fmt.Errorf("max file size must be greater than 0")
	}

	return nil
}

// GetMaxFileSizeBytes возвращает максимальный размер файла в байтах
func (c *Config) GetMaxFileSizeBytes() int64 {
	return c.Upload.MaxFileSize * 1024 * 1024 // MB to bytes
}

// GetAddress возвращает адрес для запуска сервера
func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// SaveExample создает пример конфигурационного файла
func (c *Config) SaveExample(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Определяем формат по расширению файла
	ext := filepath.Ext(path)
	switch strings.ToLower(ext) {
	case ".yaml", ".yml":
		encoder := yaml.NewEncoder(file)
		encoder.SetIndent(2)
		return encoder.Encode(c)
	case ".json":
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		return encoder.Encode(c)
	default:
		// По умолчанию JSON
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		return encoder.Encode(c)
	}
}
