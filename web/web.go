package web

import (
	"github.com/polis-mail-ru-golang-1/t2-invert-index-search-RGrouse/controller"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type Web struct {
	address 	string
	controller 	controller.Controller
}

func New(listen string, c controller.Controller) Web {
	return Web{
		address:listen,
		controller:c,
	}
}

func (web *Web) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/search", web.controller.SearchResultHandler)
	mux.HandleFunc("/add", web.controller.AddToIndexHandler)
	mux.HandleFunc("/main", web.controller.SearchHandler)
	mux.HandleFunc("/", web.controller.DefaultHandler)
	mux.Handle("/res/", http.StripPrefix(
		"/res/",
		http.FileServer(http.Dir("./view/static")),
	))

	server := http.Server{
		Addr:         web.address,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Info().Msgf("Старт сервера на интерфейсе %v", web.address)
	return server.ListenAndServe()
}
