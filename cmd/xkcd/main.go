package main

import (
	"flag"
	"log"
	"myapp/internal/config"
	"myapp/internal/interrupt"
	"myapp/internal/search"
	"myapp/pkg/index"
)

// структура для хранения конфигурации

func ParseFlags() (string, string, bool) {
	pathConfig := flag.String("c", "config.yaml", "path to the configuration file")
	searchLine := flag.String("s", "", "line for search comics by database.json")
	isIndexSearch := flag.Bool("i", false, "enables index search")
	flag.Parse()

	return *pathConfig, *searchLine, *isIndexSearch
}

func main() {
	// Единственный флаг
	pathConfig, searchLine, isIndexSearch := ParseFlags()

	// Работа с конфигом
	config, err := config.LoadConfig(pathConfig)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Создаем бд
	db, err := interrupt.CreateDatabase(config.SourceURL, config.DbFile, config.Parallel)
	if err != nil {
		log.Fatal("error create database: ")
	}
	defer db.Close()

	index, err := index.CreateIndex(config.IndexFile, db)
	if err != nil {
		log.Fatal("error create index: " + err.Error())
	}

	// выполняем поиск
	search.Search(searchLine, db, index, isIndexSearch, true)

}
