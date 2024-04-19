package database

import (
	"encoding/json"
	"os"
	"sync"
)

// ComicData структура для хранения информации о каждом комиксе
type ComicData struct {
	Img      string   `json:"url"`
	Keywords []string `json:"keywords"`
}

// структура для хранения базы данных комиксов
type ComicsDatabase struct {
	Path    string
	Mutex   sync.Mutex
	Records map[int]ComicData
}

// открывает базу данных, читает её содержимое из файла или создаёт новый файл, если он не существует
func Open(path string) (*ComicsDatabase, error) {
	db := ComicsDatabase{
		Path: path,
	}

	if _, err := os.Stat(db.Path); os.IsNotExist(err) {
		db.Records = make(map[int]ComicData)
		bytes, err := json.MarshalIndent(db.Records, "", "\t")
		if err != nil {
			return nil, err
		}

		err = os.WriteFile(db.Path, bytes, 0644)
		if err != nil {
			return nil, err
		}
	} else if err == nil {
		bytes, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(bytes, &db.Records)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	return &db, nil
}

// сохраняет данные базы данных в файл
func (db *ComicsDatabase) Close() error {
	bytes, err := json.MarshalIndent(db.Records, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(db.Path, bytes, 0644)
}

// чтобы добавлять каждый успешный комикс

// добавляет новый комикс в базу данных
func (db *ComicsDatabase) AddComic(id int, comic ComicData) {
	db.Mutex.Lock()
	defer db.Mutex.Unlock()
	db.Records[id] = comic
}
