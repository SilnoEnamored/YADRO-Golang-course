package xkcd

import (
	"encoding/json"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"net/http"
	"net/url"
	"strconv"
)

type Comics struct {
	ID         int    `json:"num"`
	ImageURL   string `json:"img"`
	Transcript string `json:"transcript"`
	AltText    string `json:"alt"`
}

const latestComicEndpoint = "info.0.json"

// Загружает данные комикса по URL
func fetchComic(comicURL string) (Comics, error) {
	// Запрос
	resp, err := http.Get(comicURL)
	if err != nil {
		return Comics{}, err
	}
	defer resp.Body.Close()

	// Проверка на корректность ответа
	if resp.StatusCode != http.StatusOK {
		return Comics{}, fmt.Errorf("failed to fetch comic: HTTP %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	// Читаем Json из запроса и преобразуем в структуру comic
	var comic Comics
	if err := json.NewDecoder(resp.Body).Decode(&comic); err != nil {
		return Comics{}, err
	}

	return comic, nil
}

// Данные последнего комикса для определения общего количества
func fetchLatestComic(baseURL string) (Comics, error) {
	urlLastComic, err := url.JoinPath(baseURL, latestComicEndpoint)
	if err != nil {
		return Comics{}, err
	}

	return fetchComic(urlLastComic)
}

// Загружает данные всех комиксов XKCD
func FetchAllComics(baseURL string) ([]Comics, error) {
	latestComic, err := fetchLatestComic(baseURL)
	if err != nil {
		return nil, err
	}
	totalComics := latestComic.ID // Используем ID последнего комикса как общее количество комиксов

	// Прогрессбар
	bar := progressbar.Default(int64(totalComics))

	comicsSlice := make([]Comics, 0, totalComics)

	for i := 1; i <= totalComics; i++ {
		bar.Add(1)
		comicURL, err := url.JoinPath(baseURL, strconv.Itoa(i), latestComicEndpoint)
		if err != nil {
			continue // Пропускаем комиксы с ошибками в URL
		}

		comic, err := fetchComic(comicURL)
		if err != nil {
			continue // Пропускаем комиксы с ошибками при загрузке
		}

		comicsSlice = append(comicsSlice, comic)
	}

	return comicsSlice, nil
}
