package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"math"
	"myapp/pkg/database"
	"myapp/pkg/words"
	"myapp/pkg/xkcd"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Config структура для хранения конфигурации
type Config struct {
	SourceURL string `yaml:"source_url"`
	DbFile    string `yaml:"db_file"`
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

	return config, nil
}

// Обрабатывает флаги
func parseFlags() (bool, int) {
	showDb := flag.Bool("o", false, "Display database content")
	numComics := flag.Int("n", math.MaxInt, "Number of comics to fetch")

	flag.Parse()
	return *showDb, *numComics
}

func rootDir() string {
	// Возвращает директорию, в которой расположен исполняемый файл
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func main() {
	showDb, numComics := parseFlags()
	config, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	dbPath := filepath.Join(rootDir(), config.DbFile)
	db := database.ComicsDatabase{}

	// Проверка на существование БД
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		fmt.Println(config.DbFile + " does not exist, creating new one.")
	} else {
		fmt.Println(config.DbFile + " exists, loading existing database.")
		// Загрузка существующей БД
		err = db.Load(dbPath)
		if err != nil {
			log.Fatalf("Error loading database: %v", err)
		}
	}

	if showDb {
		// Проверка на количество
		for i := 1; i <= numComics; i++ {
			if i > len(db)+1 {
				return
			}
			data := db[strconv.Itoa(i)]

			// Выводит на экран содержимое БД
			fmt.Printf("ID: %v\nImage: %s\nKeywords: %v\n", i, data.Img, data.Keywords)
		}
		return
	}

	// Загрузка всех комиксов
	comicsSlice, err := xkcd.FetchAllComics(config.SourceURL)
	if err != nil {
		log.Fatalf("Error fetching all comics: %v", err)
	}

	// Обработка и сохранение комиксов
	for _, comic := range comicsSlice {
		// Комбинируем транскрипт и alt текст и нормализуем ключевые слова
		text := comic.Transcript + " " + comic.AltText
		keywords, err := words.Normalization(strings.Fields(text))
		if err != nil {
			log.Fatalf("Error Normalization: %v, comic: %v", err, comic)
			return
		}

		// Сохраняем данные о комиксе в структуру БД
		db[fmt.Sprintf("%d", comic.ID)] = database.ComicData{
			Img:      comic.ImageURL,
			Keywords: keywords,
		}
	}

	// Сохраняем базу данных в файл
	err = db.Save(dbPath)
	if err != nil {
		log.Fatalf("Error saving database: %v", err)
	}
}
