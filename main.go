package main

import (
	"bufio"
	"github.com/go-pg/pg"
	"github.com/polis-mail-ru-golang-1/t2-invert-index-search-RGrouse/config"
	"github.com/polis-mail-ru-golang-1/t2-invert-index-search-RGrouse/controller"
	"github.com/polis-mail-ru-golang-1/t2-invert-index-search-RGrouse/model/db_model"
	"github.com/polis-mail-ru-golang-1/t2-invert-index-search-RGrouse/model/interfaces"
	"github.com/polis-mail-ru-golang-1/t2-invert-index-search-RGrouse/model/map_model"
	"github.com/polis-mail-ru-golang-1/t2-invert-index-search-RGrouse/view"
	"github.com/polis-mail-ru-golang-1/t2-invert-index-search-RGrouse/web"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"sync"
)

func main() {
	cfg, err := config.Load()
	die(err)

	level, err := zerolog.ParseLevel(cfg.LogLevel)
	die(err)
	zerolog.MessageFieldName = "msg"
	log.Level(level)

	log.Print(cfg)

	var m interfaces.InvertedIndexModel

	switch cfg.ModelType {
	case "MAP":
		m = map_model.New()
	case "DB":
		pgOpt, err := pg.ParseURL(cfg.PgSQL)
		die(err)
		pgdb := pg.Connect(pgOpt)
		m = db_model.New(pgdb)
	default:
		panic("Неправильный параметр MODEL")
	}

	indexFilesInFolder(cfg.SearchingFolder, m)

	v, err := view.New()
	die(err)

	c := controller.New(v, m)

	webServer := web.New(cfg.Listen, c)
	die(webServer.Start())
}

func indexFilesInFolder(searchingfolder string, iim interfaces.InvertedIndexModel) {
	log.Info().Msg("Индексируем файлы из директории "+searchingfolder)

	filesInfos, err := ioutil.ReadDir(searchingfolder)
	die(err)

	type fileWordsEntry struct {
		Source string
		CountedWords map[string]int
	}

	var filesCountInFolder int
	for _, fileInfo := range filesInfos {
		if (!fileInfo.IsDir()) {
			filesCountInFolder++
		}
	}

	ch := make(chan fileWordsEntry, filesCountInFolder)

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		for i:=0; i< filesCountInFolder; i++ {
			fileWordsEntry := <-ch
			err := iim.AttachWeightedWords(fileWordsEntry.Source, fileWordsEntry.CountedWords) //слушаем канал и добавляем в общий индекс посчитанные слова
			if(err!=nil) {
				log.Error().Msgf("Ошибка при добавлении слов в индекс %v", err)
				return
			}
		}
		wg.Done()
	}()

	for _, fileInfo := range filesInfos {
		wg.Add(1)

		go func(f os.FileInfo) {
			defer wg.Done()

			if (!f.IsDir()) {
				words, err := wordsInFile(searchingfolder + "/" + f.Name()) //разбиваем файл по словам
				if(err!=nil) {
					log.Error().Msgf("Ошибка при обработке файла %v", err)
					return
				}
				weighted := interfaces.CountWords(words)      //считаем, сколько раз слово появилось в файле

				ch <- fileWordsEntry{f.Name(), weighted} //пишем в канал источник и карту посчитанных слов
			}
		}(fileInfo)
	}
	wg.Wait()	//ждем пока все файлы проиндексируются
	log.Info().Msg("Закончили индексирование")
}

func wordsInFile(path string) ([]string, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	words := []string{}
	for scanner.Scan() {
		cleanedWords := interfaces.WordsInString(scanner.Text())
		words = append(words, cleanedWords...)
	}

	for i, _ := range words {
		words[i]=interfaces.StemWord(words[i])
	}

	return words, nil
}

func die(err error) {
	if err != nil {
		panic(err)
	}
}