package main

import (
	"flag"
	"fmt"
	"github.com/kljensen/snowball"
	"strings"
)

func main() {
	sFlag := flag.String("s", "", "Input string to be normalized")
	flag.Parse()

	var words []string

	if *sFlag != "" {
		*sFlag = strings.ReplaceAll(*sFlag, "i'll", "i will")
		words = strings.Fields(*sFlag)
	} else {
		// Если флаг -s не использовался, парсим аргументы как список слов
		words = flag.Args()
	}

	result := normalization(words)
	fmt.Println(strings.Join(result, " "))
}

func normalization(words []string) []string {

	ignoredWords := map[string]bool{
		"of": true, "a": true, "the": true, "will": true, "i": true, "you": true, "as": true, "are": true, "me": true,
		// Можно добавить другие слова, которые нужно игнорировать
	}

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
