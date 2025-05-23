package http

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"quotes/internal/models"
	"quotes/internal/service"
	errs "quotes/pkg/errors"
	"strconv"
)

type Handler struct {
	service *service.QuotesService
}

func NewHandler(srv *service.QuotesService) (*Handler, error) {
	return &Handler{
		service: srv,
	}, nil
}

func (h *Handler) CreateQuote(w http.ResponseWriter, r *http.Request) {
	var quote models.Quote
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&quote); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	quoteId, err := h.service.CreateQuote(&quote)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidQuoteField) {
			http.Error(w, "invalid quote field", http.StatusUnprocessableEntity)
			return
		}
		if errors.Is(err, errs.ErrInvalidAuthorField) {
			http.Error(w, "invalid author field", http.StatusUnprocessableEntity)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(map[string]uint64{"id": quoteId})
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (h *Handler) GetQuotes(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	var author *string
	if authorParam := query.Get("author"); authorParam != "" {
		author = &authorParam
	}
	
	quotes, err := h.service.GetQuotes(author)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(quotes)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *Handler) GetRandomQuote(w http.ResponseWriter, r *http.Request) {
	quote, err := h.service.GetRandomQuote()
	if err != nil {
		if errors.Is(err, errs.ErrQuoteNotFound) {
			http.Error(w, "quote not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(quote)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (h *Handler) DeleteQuote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteQuote(id)
	if err != nil {
		if errors.Is(err, errs.ErrQuoteNotFound) {
			http.Error(w, "quote not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, errs.ErrInvalidQuoteId) {
			http.Error(w, "invalid id", http.StatusUnprocessableEntity)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
