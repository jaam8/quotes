package storage

import (
	"errors"
	"quotes/internal/models"
	errs "quotes/pkg/errors"
	"sync"
)

type DBAdapter struct {
	idCounter uint64
	byId      map[uint64]*models.Quote
	byAuthor  map[string][]*models.Quote
	mu        sync.RWMutex
}

func NewDBAdapter() *DBAdapter {
	return &DBAdapter{
		idCounter: 0,
		byId:      make(map[uint64]*models.Quote),
		byAuthor:  make(map[string][]*models.Quote),
		mu:        sync.RWMutex{},
	}
}

func (db *DBAdapter) SaveQuote(quote *models.Quote) (uint64, error) {
	if quote == nil {
		return 0, errors.New("quote is nil")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	db.idCounter += 1
	quote.Id = db.idCounter

	db.byId[db.idCounter] = quote
	db.byAuthor[quote.Author] = append(db.byAuthor[quote.Author], quote)

	return quote.Id, nil
}

func (db *DBAdapter) GetQuote(id uint64) (*models.Quote, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	quote, ok := db.byId[id]
	if !ok {
		return nil, errs.ErrQuoteNotFound
	}

	return quote, nil
}

func (db *DBAdapter) GetQuotes(author *string) ([]*models.Quote, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if author == nil {
		var quotes []*models.Quote
		for _, q := range db.byId {
			quotes = append(quotes, q)
		}
		return quotes, nil
	}

	quotes, ok := db.byAuthor[*author]
	if !ok {
		return []*models.Quote{}, nil
	}

	return quotes, nil
}

func (db *DBAdapter) DeleteQuote(id uint64) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	quote, ok := db.byId[id]
	if !ok {
		return errs.ErrQuoteNotFound
	}
	delete(db.byId, id)

	quotes := db.byAuthor[quote.Author]
	for i, q := range db.byAuthor[quote.Author] {
		if q.Id == quote.Id {
			quotes = append(quotes[:i], quotes[i+1:]...)
			break
		}
	}

	if len(quotes) == 0 {
		delete(db.byAuthor, quote.Author)
	} else {
		db.byAuthor[quote.Author] = quotes
	}

	return nil
}
