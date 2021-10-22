package handler

import (
	"avilego.me/recent_news/factory"
	"avilego.me/recent_news/handler/api"
	"net/http"
)

func NewServerHttpHandler() http.Handler {
	mux := http.NewServeMux()
	configRoutes(mux)
	return mux
}

func configRoutes(mux *http.ServeMux) {
	mux.Handle("/api/search", api.ApiSearchHandler{Finder: factory.Finder()})
}