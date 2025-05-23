package http

import "net/http"

type Handler struct {
	//	service
}

func NewHandler() *Handler {
	var handler Handler
	// handler.service = service
	return &handler
}

func (h *Handler) CreateQuote(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) GetQuotes(w http.ResponseWriter, r *http.Request) {
	//	+ query params
}

func (h *Handler) GetRandomQuote(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) DeleteQuote(w http.ResponseWriter, r *http.Request) {
	//	+ path params
}
