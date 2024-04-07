package database

import (
	"encoding/json"
	"io/ioutil"
)

// ComicData структура для хранения информации о каждом комиксе.
type ComicData struct {
	Img      string   `json:"url"`
	Keywords []string `json:"keywords"`
}

// ComicsDatabase структура для хранения базы данных комиксов.
type ComicsDatabase map[string]ComicData

// Save записывает базу данных комиксов в файл JSON.
func (db ComicsDatabase) Save(filename string) error {
	data, err := json.MarshalIndent(db, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

// Load читает базу данных комиксов из файла JSON.
func (db *ComicsDatabase) Load(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, db)
}
