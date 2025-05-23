package storage

import (
	"quotes/internal/models"
	"sync"
)

type DBAdapter struct {
	db map[uint64]*models.Quote
	mu sync.Mutex
}

func NewDBAdapter() *DBAdapter {
	return &DBAdapter{
		db: make(map[uint64]*models.Quote),
		mu: sync.Mutex{},
	}
}

func (db *DBAdapter) CreateQuote(quote *models.Quote) (uint64, error) {
	return 0, nil
}

func (db *DBAdapter) GetQuote(id uint64) (*models.Quote, error) {
	return nil, nil
}

func (db *DBAdapter) GetQuotes(author *string) ([]*models.Quote, error) {
	return nil, nil
}

func (db *DBAdapter) DeleteQuote(id uint64) (*models.Quote, error) {
	return nil, nil
}
