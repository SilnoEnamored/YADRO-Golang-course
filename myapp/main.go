package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/kljensen/snowball"
	"os"
	"strings"
)

func main() {
	sFlag := flag.String("s", "", "Input string to be normalized")
	flag.Parse()

	var words []string
	words = strings.Fields(*sFlag)

	result := normalization(words)
	fmt.Println(strings.Join(result, " "))
}

func normalization(words []string) []string {

	//Файл со стоп словами(english), которые необходимо игнорировать
	file, err := os.Open("stop_words_english.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	//Мапа для стоп слов
	var ignoredWords = make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ignoredWords[scanner.Text()] = true
	}

	//Мапа для повторяющихся слов
	repeated := make(map[string]struct{})

	var result []string
	for _, word := range words {
		word = strings.ToLower(word)
		if !ignoredWords[word] {
			stemmed, err := snowball.Stem(word, "english", true)
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
