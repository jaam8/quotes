package ports

import "quotes/internal/models"

type StorageAdapter interface {
	CreateQuote(quote *models.Quote) (uint64, error)
	GetQuote(id uint64) (*models.Quote, error)
	GetQuotes(author *string) ([]*models.Quote, error)
	DeleteQuote(id uint64) (*models.Quote, error)
}
