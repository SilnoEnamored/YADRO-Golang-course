package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/kljensen/snowball"
	"os"
	"strings"
)

var ignoredWords = make(map[string]bool)

func main() {
	sFlag := flag.String("s", "", "Input string to be normalized")
	flag.Parse()

	loadStopWords("stop_words_english.txt")

	var words []string
	words = strings.Fields(*sFlag)

	result := normalization(words, ignoredWords)
	fmt.Println(strings.Join(result, " "))
}

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

func loadStopWords(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ignoredWords[scanner.Text()] = true
	}
}

func normalization(words []string, ignoredWords map[string]bool) []string {

	//Мапа для повторяющихся слов
	repeated := make(map[string]struct{})

	var result []string
	for _, word := range words {
		word = strings.ToLower(word)
		if !ignoredWords[word] {
			stemmed, err := snowball.Stem(deleteUnnecessary(word), "english", true)
			if err != nil {
				fmt.Println("Error stemming word:", err)
				continue
			}
			if _, ok := repeated[stemmed]; !ok {
				result = append(result, stemmed)
				repeated[stemmed] = struct{}{}
			}
		}
	}
	return result
}
