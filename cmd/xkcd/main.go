package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"myapp/internal/interrupt"
)

// структура для хранения конфигурации
type Config struct {
	SourceURL string `yaml:"source_url"`
	DbFile    string `yaml:"db_file"`
	Parallel  int    `yaml:"parallel"`
}

// Загружает конфигурацию из YAML файла
func loadConfig(path string) (Config, error) {
	var config Config

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("unable to load config file: %w", err)
	}

	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return Config{}, fmt.Errorf("yaml decode error: %w", err)
	}

	// Проверяем, что parallel больше 0
	if config.Parallel <= 0 {
		config.Parallel = 1 // устанавливаем по умолчанию в 1, если значение не задано или отрицательно
	}

	return config, nil
}

func main() {
	// Единственный флаг
	configPath := flag.String("c", "config.yaml", "Path to the configuration file")
	flag.Parse()

	// Работа с конфигом
	config, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Создаем бд
	db, err := interrupt.CreateDatabase(config.SourceURL, config.DbFile, config.Parallel)
	if err != nil {
		log.Fatal("error create database: ")
	}
	defer db.Close()

}
