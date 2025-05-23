package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	h "quotes/internal/http"
)

func main() {
	r := mux.NewRouter()

	handler := h.NewHandler()

	r.HandleFunc("/quotes", handler.CreateQuote).Methods("POST")
	r.HandleFunc("/quotes", handler.GetQuotes).Methods("GET")
	r.HandleFunc("/quotes/random", handler.GetRandomQuote).Methods("PUT")
	r.HandleFunc("/quotes/:id", handler.DeleteQuote).Methods("DELETE")

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
