package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/search", app.searchHome)
	mux.HandleFunc("/search/", app.searchWord)
	mux.HandleFunc("/query", app.searchHome)
	mux.HandleFunc("/query/", app.searchWord)
	return mux
}
