package controller

import (
	"fmt"
	"github.com/polis-mail-ru-golang-1/t2-invert-index-search-RGrouse/model/interfaces"
	"github.com/polis-mail-ru-golang-1/t2-invert-index-search-RGrouse/view"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

type Controller struct {
	view  view.View
	model interfaces.InvertedIndexModel
}

func New(v view.View, m interfaces.InvertedIndexModel) Controller {
	return Controller{
		view:  v,
		model: m,
	}
}

func (c Controller) SearchHandler(w http.ResponseWriter, r *http.Request) {
	c.checkTemplateExec(c.view.Search(w), w, r)
}
func (c Controller) SearchResultHandler(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")
	if query != "" {
		q := strings.ToLower(query)

		result, err := c.model.SearchByString(q)
		if err!=nil{
			c.error(w, r, "Ошибка при поиске", 500)
			return
		}
		c.checkTemplateExec(c.view.SearchResult(q, result, w), w, r)
	} else {
		http.Redirect(w, r, "/main", 301)
	}
}
func (c Controller) DefaultHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/main", 301)
}

func (c Controller) AddToIndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET"{
		c.checkTemplateExec(c.view.AddToIndex(w), w, r)
	} else if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			log.Error().Msgf("Ошибка при обработке формы %v", err)
			c.error(w,r,"Ошибка на сервере",500)
			return
		}
		text := r.FormValue("text")
		source := r.FormValue("source")
		if (text!="" && source!=""){
			strs := strings.Split(text, " ")
			weighted := interfaces.CountWords(strs)
			err := c.model.AttachWeightedWords(source, weighted)
			if err!=nil{
				c.error(w,r,"Ошибка при добавлении", 500)
				return
			}
			c.checkTemplateExec(c.view.AddToIndexPopup(true, w), w, r)
		} else {
			//log.Error().Msgf("Неподходящая форма от пользователя")
			//c.error(w,r,"Неподходящая форма",400)
			c.checkTemplateExec(c.view.AddToIndexPopup(false, w), w, r)
		}
	}
}

func (c Controller) error(w http.ResponseWriter, r *http.Request, err string, status int) {
	log.Info().Msgf("Ошибка %s код состояния %d", err, status)
	if err := c.view.Error(err, status, w); err != nil {
		log.Error().Err(err).Msgf("Ошибка при обработке шаблона ошибки %s", err)
		fmt.Fprintln(w, "Ошибка на сервере")
	}
}

func (c Controller) checkTemplateExec(err error, w http.ResponseWriter, r *http.Request) {
	if err != nil {
		log.Error().Msgf("Ошибка при обработке шаблона %s", err)
		c.error(w, r, "Ошибка на сервере", 500)
	}
}