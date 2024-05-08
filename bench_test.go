package bench_test

import (
	"log"
	"testing"

	"myapp/internal/config"
	"myapp/internal/interrupt"
	"myapp/internal/search"
	"myapp/pkg/index"
)

func BenchmarkDefaultSearch(b *testing.B) {
	config, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal("Error loading config: " + err.Error())
	}

	db, err := interrupt.CreateDatabase(config.SourceURL, config.DbFile, config.Parallel)
	if err != nil {
		log.Fatal("Error create database: " + err.Error())
	}
	defer db.Close()

	index, err := index.CreateIndex(config.IndexFile, db)
	if err != nil {
		log.Fatal("Error create index: " + err.Error())
	}

	b.ReportAllocs() //Замеряю только поиск, без загрузки конфигурации и создания бд
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		search.Search("I'm following your questions", db, index, false, false)
		b.StopTimer()
	}
}

func BenchmarkIndexSearch(b *testing.B) {
	config, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal("Error loading config: " + err.Error())
	}

	db, err := interrupt.CreateDatabase(config.SourceURL, config.DbFile, config.Parallel)
	if err != nil {
		log.Fatal("Error create database: " + err.Error())
	}
	defer db.Close()

	index, err := index.CreateIndex(config.IndexFile, db)
	if err != nil {
		log.Fatal("Error create index: " + err.Error())
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		search.Search("I'm following your questions", db, index, true, false)
		b.StopTimer()
	}
}
