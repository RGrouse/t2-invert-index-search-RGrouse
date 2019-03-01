package map_model

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

var wordsEntry1, wordsEntry2, wordsEntry3 WordsEntry
func init() {
	wordsEntry1 = WordsEntry{
		"one",
		map[string]int{
			"1": 1,
			"2": 2,
			"3": 3,
		},
	}
	wordsEntry2 = WordsEntry{
		"two",
		map[string]int{
			"1": 1,
			"2": 2,
			"4": 4,
			"3": 3,
		},
	}
	wordsEntry3 = WordsEntry{}
}

func TestNewIndex(t *testing.T){
	m1 := gindexmap
	m1["1"]=nil
	NewIndexMap()
	m2 := gindexmap

	if reflect.DeepEqual(m1, m2) {
		t.Errorf("%v is equal to unexpected %v", m1, m2)
	}
}

func TestCountWords(t *testing.T) {
	s := []string{"1", "3", "2", "2", "3", "3"}
	m := CountWords(s)
	actm := map[string]int{
		"1": 1,
		"2": 2,
		"3": 3,
	}
	if !reflect.DeepEqual(*m, actm) {
		t.Errorf("%v is not equal to expected %v", m, actm)
	}

	emptym := CountWords([]string{})
	if len(*emptym) != 0 {
		t.Errorf("%v is not empty (0 len arg)", emptym)
	}

	emptym = CountWords(nil)
	if len(*emptym) != 0 {
		t.Errorf("%v is not empty (nil arg)", emptym)
	}
}

func TestAttachCountedWordsFromChannel(t *testing.T) {
	NewIndexMap()

	resultMap1 := map[string][]IndexEntry{
		"1": []IndexEntry{IndexEntry{"one", 1}},
		"2": []IndexEntry{IndexEntry{"one", 2}},
		"3": []IndexEntry{IndexEntry{"one", 3}},
	}
	resultMap2 := map[string][]IndexEntry{
		"1": []IndexEntry{IndexEntry{"one", 1}, IndexEntry{"two", 1}},
		"2": []IndexEntry{IndexEntry{"one", 2}, IndexEntry{"two", 2}},
		"3": []IndexEntry{IndexEntry{"one", 3}, IndexEntry{"two", 3}},
		"4": []IndexEntry{IndexEntry{"two", 4}},
	}

	ch := make(chan WordsEntry)

	go func() {
		AttachCountedWordsFromChannel(ch, 3)
	}()

	ch <- wordsEntry1
	time.Sleep(100000)
	if !reflect.DeepEqual(gindexmap, resultMap1) {
		t.Errorf("%v is not equal to expected %v (attach first map)", gindexmap, resultMap1)
	}

	ch <- wordsEntry2
	time.Sleep(100000)
	if !reflect.DeepEqual(gindexmap, resultMap2) {
		t.Errorf("%v is not equal to expected %v (attach second map)", gindexmap, resultMap2)
	}

	ch <- wordsEntry3
	time.Sleep(100000)
	if !reflect.DeepEqual(gindexmap, resultMap2) {
		t.Errorf("%v is not equal to expected %v (attach empty map)", gindexmap, resultMap2)
	}
}
func TestSearchByWords(t *testing.T){
	NewIndexMap()

	ch := make(chan WordsEntry)
	go func() {
		AttachCountedWordsFromChannel(ch, 2)
	}()
	ch <- wordsEntry1
	ch <- wordsEntry2
	time.Sleep(100000)

	m := SearchByWords([]string{})
	if len(m)!=0{
		t.Errorf("%v is not empty (0 len arg)", m)
	}

	m = SearchByWords(nil)
	if len(m)!=0{
		t.Errorf("%v is not empty (nil arg)", m)
	}

	m = SearchByWords([]string{"1"})
	resultSearch1 := map[string]int{
		"one": 1,
		"two": 1,
	}
	if !reflect.DeepEqual(m, resultSearch1) {
		t.Errorf("%v is not equal to expected %v", m, resultSearch1)
	}

	m = SearchByWords([]string{"1", "4"})
	resultSearch2 := map[string]int{
		"one": 1,
		"two": 5,
	}
	if !reflect.DeepEqual(m, resultSearch2) {
		t.Errorf("%v is not equal to expected %v", m, resultSearch2)
	}
}

func TestSearchByString(t *testing.T){
	NewIndexMap()

	ch := make(chan WordsEntry)
	go func() {
		AttachCountedWordsFromChannel(ch, 2)
	}()
	ch <- wordsEntry1
	ch <- wordsEntry2
	time.Sleep(100000)

	m1 := SearchByString("")
	m2 := SearchByWords([]string{})
	if !reflect.DeepEqual(m1, m2){
		t.Errorf("%v is not equal to expected %v (0 len arg)", m1, m2)
	}

	m1 = SearchByString("1")
	m2 = SearchByWords([]string{"1"})
	if !reflect.DeepEqual(m1, m2) {
		t.Errorf("%v is not equal to expected %v", m1, m2)
	}

	m1 = SearchByString("1 4")
	m2 = SearchByWords([]string{"1", "4"})
	if !reflect.DeepEqual(m1, m2) {
		t.Errorf("%v is not equal to expected %v", m1, m2)
	}
}

func TestSortAndPrintResult(t *testing.T) {
	unsorted := map[string]int{
		"ex1.txt":1,
		"ex6.txt":6,
		"ex3.txt":3,
		"ex2.txt":2,
		"ex5.txt":5,
		"ex4.txt":4,
	}
	sortedResultString := "- ex6.txt; совпадений - 6\n"+
		"- ex5.txt; совпадений - 5\n"+
		"- ex4.txt; совпадений - 4\n"+
		"- ex3.txt; совпадений - 3\n"+
		"- ex2.txt; совпадений - 2\n"+
		"- ex1.txt; совпадений - 1\n"

	buf := new(bytes.Buffer)
	SortAndPrintResult(unsorted, buf)
	if buf.String()!=sortedResultString {
		t.Errorf("%v is not equal to expected %v", buf.String(), sortedResultString)
	}
}