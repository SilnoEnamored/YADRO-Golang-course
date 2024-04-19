package xkcd

import (
	"encoding/json"
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

// Клиент для работы с Api
type ComicClient struct {
	BaseURL   string
	Endpoints map[string]string
}

func NewComicClient(baseURL string) *ComicClient {
	endpoints := map[string]string{
		"FetchComic": "info.0.json",
	}
	return &ComicClient{
		BaseURL:   baseURL,
		Endpoints: endpoints,
	}
}

// находит ID последнего доступного комикса, сложность O(logN)
func (client *ComicClient) FetchLatestComicID() int {
	id := 1
	for step := 10; ; step *= 10 {
		if client.FetchComic(id).ID == 0 {
			break
		}
		id *= 10
	}

	low, high := id/10, id
	for low < high {
		mid := low + (high-low)/2
		if client.FetchComic(mid).ID == 0 {
			high = mid
		} else {
			low = mid + 1
		}
	}

	return low - 1
}

// возвращает комикс по его ID
func (client *ComicClient) FetchComic(comicID int) Comics {
	url := client.buildComicURL(comicID)
	return client.fetchComicData(url)
}

// извлекает данные о комиксе по URL
func (client *ComicClient) fetchComicData(endpoint string) Comics {
	response, err := http.Get(endpoint)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	var comic Comics
	if response.StatusCode == http.StatusOK {
		err := json.NewDecoder(response.Body).Decode(&comic)
		if err != nil {
			panic(err)
		}
	}

	return comic
}

// создает URL для получения комикса по ID
func (client *ComicClient) buildComicURL(comicID int) string {
	comicURL, err := url.JoinPath(client.BaseURL, strconv.Itoa(comicID), client.Endpoints["FetchComic"])
	if err != nil {
		panic(err)
	}

	return comicURL
}
