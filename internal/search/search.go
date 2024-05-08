package search

import (
	"log"
	"sort"
	"strings"

	"myapp/pkg/database"
	"myapp/pkg/index"
	"myapp/pkg/words"
)

// SearchResult хранит пару: ID комикса и его оценку поисковой релевантности
type SearchResult struct {
	Id    int
	Score int // Оценка соответсвия
}

// Преобразует строку запроса в множество нормализованных слов
func prepareSearchTerms(request string) (map[string]bool, error) {
	wordsRequest := strings.Fields(request)
	searchTerms := make(map[string]bool)

	normalizedWords, err := words.Normalization(wordsRequest)
	if err != nil {
		log.Printf("Error during normalization: %v", err)
		return nil, err
	}

	for _, word := range normalizedWords {
		searchTerms[word] = true
	}

	return searchTerms, nil
}

// Вычисляет оценку релевантности для каждого ID
func calculateScores(searchTerms map[string]bool, getIDs func(string) []int) map[int]int {
	scores := make(map[int]int)
	for term := range searchTerms {
		for _, id := range getIDs(term) {
			scores[id]++
		}
	}
	return scores
}

// Вычисляет релевантность с использованием базы данных
func relevanceFromDatabase(searchTerms map[string]bool, db *database.ComicsDatabase) map[int]int {
	getIDs := func(term string) []int {
		var ids []int
		for id, comics := range db.Records {
			for _, keyword := range comics.Keywords {
				if keyword == term {
					ids = append(ids, id)
				}
			}
		}
		return ids
	}
	return calculateScores(searchTerms, getIDs)
}

// Вычисляет релевантность с использованием индексного файла
func relevanceFromIndex(searchTerms map[string]bool, index *index.Index) map[int]int {
	return calculateScores(searchTerms, index.GetComics)
}

// Сортирует результаты поиска по убыванию релевантности
func sortResultsByRelevance(scores map[int]int) []SearchResult {
	results := make([]SearchResult, 0, len(scores))
	for id, score := range scores {
		results = append(results, SearchResult{Id: id, Score: score})
	}
	sort.Slice(results, func(i, j int) bool { return results[i].Score > results[j].Score })
	return results
}

func init() {
	log.SetFlags(0) // Убирает все декорации при использовании log.Printf, включая время
}

// Печатает заданное количество наиболее релевантных комиксов,так же isPrint отвечает за отображение,false- для удобства при бенчмарке
func displayTopResults(results []SearchResult, count int, isPrint bool) {
	if !isPrint {
		return
	}
	for i := 0; i < minimum(count, len(results)); i++ {
		log.Printf("%d: https://xkcd.com/%d/info.0.json Score: %d\n", i+1, results[i].Id, results[i].Score)
	}
}

// Возвращает минимальное значение из двух
func minimum(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Выполняет поиск, агрегируя все шаги
func Search(request string, db *database.ComicsDatabase, index *index.Index, isIndexSearch bool, isPrint bool) {
	searchTerms, err := prepareSearchTerms(request)
	if err != nil {
		log.Printf("Failed to prepare search terms: %v", err)
		return
	}

	var scores map[int]int
	if isIndexSearch {
		scores = relevanceFromIndex(searchTerms, index)
	} else {
		scores = relevanceFromDatabase(searchTerms, db)
	}

	results := sortResultsByRelevance(scores)

	displayTopResults(results, 10, isPrint)
}
