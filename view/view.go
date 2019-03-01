package view

import (
	"github.com/polis-mail-ru-golang-1/t2-invert-index-search-RGrouse/model/interfaces"
	"github.com/rs/zerolog/log"
	"html/template"
	"io"
	"io/ioutil"
	"strings"
)

type View struct {
	searchTmpl *template.Template
	resultTmpl  *template.Template
	errorTmpl  *template.Template
	addTmpl	*template.Template
}

func New() (View, error) {
	v := View{}

	var allFiles []string
	files, err := ioutil.ReadDir("./view/templates")
	if err != nil {
		log.Error().Msg("Ошибка при чтении директории с шаблонами")
		return v, err
	}
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".tmpl") {
			allFiles = append(allFiles, "./view/templates/"+filename)
		}
	}

	templates, err := template.ParseFiles(allFiles...)
	if err != nil {
		log.Error().Msg("Ошибка при парсинге шаблонов")
		return v, err
	}

	v.searchTmpl = templates.Lookup("search")
	v.resultTmpl = templates.Lookup("result")
	v.errorTmpl = templates.Lookup("error")
	v.addTmpl = templates.Lookup("add")

	return v, nil
}

func (v View) SearchResult(query string, searchResult []interfaces.SearchResultEntry, wr io.Writer) error {
	return v.resultTmpl.Execute(wr,
		struct {
			Query  string
			Result []interfaces.SearchResultEntry
		}{
			Query:query,
			Result:searchResult,
		})
}

func (v View) Search(wr io.Writer) error {
	return v.searchTmpl.Execute(wr, nil)
}

func (v View) AddToIndex(wr io.Writer) error {
	return v.addTmpl.Execute(wr, nil)
}
func (v View) AddToIndexPopup(isAdded bool, message string, wr io.Writer) error {
	return v.addTmpl.Execute(wr,
		struct {
			IsAdded bool
			Message string
		}{
			IsAdded: isAdded,
			Message: message,
		})
}

func (v View) Error(displayErr string, status int, wr io.Writer) error {
	return v.errorTmpl.Execute(wr,
		struct {
			Title  string
			Status int
			Error  string
		}{
			Title:  "error",
			Status: status,
			Error:  displayErr,
		})
}
