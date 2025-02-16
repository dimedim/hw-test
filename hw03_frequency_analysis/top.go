package hw03frequencyanalysis

import (
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

var (
	regularExpr = `\p{L}+(-\p{L}+)*|-{2,}`
	re          = regexp.MustCompile(regularExpr)
)

func getWords(str string) []string {
	// words := strings.Fields(str)
	str = strings.ToLower(str)
	words := re.FindAllString(str, -1)
	return words
}
