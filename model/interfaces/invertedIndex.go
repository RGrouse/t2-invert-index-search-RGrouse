package interfaces

import (
	"github.com/reiver/go-porterstemmer"
	"github.com/rs/zerolog/log"
	"strings"
)

type InvertedIndexModel interface {
	AttachWeightedWords(source string, weightedWords map[string]int) error
	SearchByString(str string) ([]SearchResultEntry, error)
	SearchByWords(words []string) ([]SearchResultEntry, error)
}

type SearchResultEntry struct {
	Source 	string
	Score 	int
}

func Check(err error){
	if err!=nil {
		log.Error().Msgf("%v",err)
	}
}
func Die(err error)  {
	if err!=nil {
		log.Fatal().Msgf("%v",err)
		panic(err)
	}
}

func CountWords(words []string) map[string]int {
	m := make(map[string]int)

	for _, word := range words {
		m[word]++
	}

	return m
}

func WordsInString(str string) []string{
	return strings.FieldsFunc(strings.ToLower(str), func(r rune) bool {
		return (r < 48 || r > 57) && (r < 97 || r > 122) && r != 45
	})
}

func StemWord(str string) string {
	return porterstemmer.StemString(str)
}