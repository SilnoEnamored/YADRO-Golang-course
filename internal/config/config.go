package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	SourceURL string `yaml:"source_url"`
	DbFile    string `yaml:"db_file"`
	IndexFile string `yaml:"index_file"`
	Parallel  int    `yaml:"parallel"`
}

// Загружает конфигурацию из YAML файла
func LoadConfig(path string) (Config, error) {
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
