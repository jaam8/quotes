package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"quotes/internal/config"
	h "quotes/internal/http"
	"quotes/internal/ports/adapters/storage"
	"quotes/internal/service"
)

func main() {

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	DBAdapter := storage.NewDBAdapter()
	srv, err := service.NewQuotesService(DBAdapter)
	handler, err := h.NewHandler(srv)

	r := mux.NewRouter()

	r.Use(h.LogMiddleware)

	r.HandleFunc("/quotes", handler.CreateQuote).Methods("POST")
	r.HandleFunc("/quotes", handler.GetQuotes).Methods("GET")
	r.HandleFunc("/quotes/random", handler.GetRandomQuote).Methods("GET")
	r.HandleFunc("/quotes/{id}", handler.DeleteQuote).Methods("DELETE")

	log.Printf("Server running at http://localhost:%d", cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), r))
}
