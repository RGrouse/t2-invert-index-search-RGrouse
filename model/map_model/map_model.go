package map_model

import (
	"github.com/polis-mail-ru-golang-1/t2-invert-index-search-RGrouse/model/interfaces"
	"github.com/rs/zerolog/log"
	"sort"
)

type MapModel struct {
	//карта вида: слово - массив из записей {источник, вес}
	indexmap map[string][]IndexEntry
}

func New() MapModel {
	return MapModel{ indexmap: make(map[string][]IndexEntry,0) }
}

type IndexEntry struct {
	Source string
	Weight int
}

func (ii MapModel) AttachWeightedWords(src string, weightedWords map[string]int) error {
	for str, w := range weightedWords {
		entries, present := ii.indexmap[str]
		if !present {
			entries = make([]IndexEntry, 0)
		}
		entries = append(entries, IndexEntry{src, w})
		ii.indexmap[str] = entries
	}
	log.Info().Msg("Добавили в индекс слова из "+ src)

	return nil
}

func (ii MapModel) SearchByString(str string) ([]interfaces.SearchResultEntry, error) {
	words := interfaces.WordsInString(str)
	return ii.SearchByWords(words)
}

func (ii MapModel) SearchByWords(words []string) ([]interfaces.SearchResultEntry, error) {
	for i, _ := range words {
		words[i] = interfaces.StemWord(words[i])
	}

	resultmap := make(map[string]int)

	for _, word := range words {
		entries, present := ii.indexmap[word]
		if !present {
			continue
		}

		for _, entry := range entries {
			resultmap[entry.Source]+=entry.Weight
		}
	}

	sortedResult := sortResult(resultmap)
	log.Info().Msgf("Поиск по словам %v, результат %v", words, sortedResult)
	return sortedResult, nil
}

func sortResult(m map[string]int) []interfaces.SearchResultEntry {
	searchResult := []interfaces.SearchResultEntry{}
	for s, w := range m {
		searchResult = append(searchResult, interfaces.SearchResultEntry{s, w})
	}

	sort.Slice(searchResult, func(i, j int) bool {
		return searchResult[i].Score > searchResult[j].Score
	})

	return searchResult
}