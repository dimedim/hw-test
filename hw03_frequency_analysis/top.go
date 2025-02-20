package hw03frequencyanalysis

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

type wordsCount struct {
	word    string
	counter int
}

func Top10(str string) []string {
	if str == "" {
		return []string{}
	}
	words := getWords(str)
	wordsMap := make(map[string]int, len(words))

	for _, val := range words {
		wordsMap[val]++
	}

	pairs := make([]wordsCount, 0, len(words))

	for key, val := range wordsMap {
		word := wordsCount{
			word:    key,
			counter: val,
		}
		pairs = append(pairs, word)
	}

	res := get10Words(pairs)

	return res
}

func get10Words(pairs []wordsCount) []string {
	res := make([]string, 0, 10)

	slices.SortFunc(pairs, func(i, j wordsCount) int {
		if j.counter == i.counter {
			return strings.Compare(i.word, j.word)
		}
		return j.counter - i.counter
	})
	j := 1
	for _, pair := range pairs {
		if j > 10 {
			break
		}
		res = append(res, pair.word)
		j++
	}

	return res
}

func getRegexp() string {
	punctSimb := `(^|\s)\p{P}{2,}($|\s)`
	words := `[^\s\p{P}]\S*[^\s\p{P}]+`
	letters := `[^\s\p{P}]`

	res := fmt.Sprintf("%s|%s|%s", punctSimb, words, letters)

	return res
}

var re = regexp.MustCompile(getRegexp())

func getWords(str string) []string {
	str = strings.ToLower(str)

	words := re.FindAllString(str, -1)
	for i, val := range words {
		words[i] = strings.TrimSpace(val)
	}

	return words
}
