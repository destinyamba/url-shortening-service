package main

import (
	"net/http"
	"url-shortening-service/internal"
	"url-shortening-service/storage/mongodb"

	"github.com/go-chi/chi/v5"
)

func main() {
	mongodb.ConnectDB()
	r := chi.NewRouter()

	r.Get("/health", internal.Health)
	r.Post("/shorten", internal.CreateShortURL)
	r.Get("/shorten/{code}", internal.GetURL)
	r.Put("/shorten/{code}", internal.UpdateURL)
	r.Delete("/shorten/{code}", internal.DeleteURL)
	r.Get("/shorten/{code}/stat", internal.GetStats)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return
	}
}
