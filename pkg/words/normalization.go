package words

import (
	"bufio"
	"github.com/kljensen/snowball"
	"os"
	"strings"
)

// Функция удаляет все в строке, что не является английской буквой
func deleteUnnecessary(s string) string {
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			result.WriteRune(r) //Можно использовать WriteByte и будет даже эффективнее по памяти и скороскти, если мы уверенны, что строка будет состоять только из ASCII символов
		}
	}
	return result.String()
}

func loadStopWords(filePath string) (map[string]bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ignoredWords := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ignoredWords[scanner.Text()] = true
	}

	return ignoredWords, nil
}

// Normalization обрабатывает список слов и возвращает нормализованный список
func Normalization(words []string) ([]string, error) {
	ignoredWords, err := loadStopWords("./pkg/words/stop_words_english.txt")
	if err != nil {
		return nil, err
	}

	repeated := make(map[string]bool)
	var result []string

	for _, word := range words {
		word = strings.ToLower(word)
		if _, exists := ignoredWords[word]; exists {
			continue
		}

		stemmed, err := snowball.Stem(deleteUnnecessary(word), "english", true)
		if err != nil {
			continue
		}

		if stemmed == "" || stemmed == "alt" || ignoredWords[stemmed] || repeated[stemmed] {
			continue
		}

		result = append(result, stemmed)
		repeated[stemmed] = true
	}

	return result, nil
}
