package handlers

import "github.com/go-file-storage/config"

// Глобальная переменная конфигурации
var appConfig *config.Config

// Init инициализирует handlers с конфигурацией
func Init(cfg *config.Config) {
	appConfig = cfg
}
