package index

import (
	"encoding/json"
	"os"

	"myapp/pkg/database"
)

// Структура для хранения индексного файла
type Index struct {
	Path           string           // Путь
	WordToComicsId map[string][]int // слова с id
}

// Создает индексный файл из бд комиксов и сохраняет его
func CreateIndex(path string, db *database.ComicsDatabase) (*Index, error) {
	index := Index{
		Path:           path,
		WordToComicsId: make(map[string][]int),
	}

	for id, comics := range db.Records {
		for _, keyWord := range comics.Keywords {
			index.AddIndex(keyWord, id)
		}
	}

	bytes, err := json.MarshalIndent(index, "", "\t")
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(index.Path, bytes, 0644)
	if err != nil {
		return nil, err
	}

	return &index, nil
}

// Добавляет id комикса к списку id
func (index *Index) AddIndex(word string, id int) {
	index.WordToComicsId[word] = append(index.WordToComicsId[word], id)
}

// Возвращает списко id комиксов, которые ассоциируются с задаными словами
func (index *Index) GetComics(word string) []int {
	return index.WordToComicsId[word]
}
