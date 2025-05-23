package service

import (
	"fmt"
	"math/rand"
	"quotes/internal/models"
	"quotes/internal/ports"
	errs "quotes/pkg/errors"
)

type QuotesService struct {
	storage ports.StorageAdapter
}

func NewQuotesService(storageAdapter ports.StorageAdapter) (*QuotesService, error) {
	return &QuotesService{
		storage: storageAdapter,
	}, nil
}

func (s *QuotesService) CreateQuote(quote *models.Quote) (uint64, error) {
	if quote.Quote == "" || quote.Quote == " " {
		return 0, errs.ErrInvalidQuoteField
	}
	if quote.Author == "" || quote.Author == " " {
		return 0, errs.ErrInvalidAuthorField
	}

	id, err := s.storage.SaveQuote(quote)
	if err != nil {
		return 0, fmt.Errorf("could not create quote: %w", err)
	}

	return id, nil
}

func (s *QuotesService) GetQuotes(author *string) ([]*models.Quote, error) {
	quotes, err := s.storage.GetQuotes(author)
	if err != nil {
		return nil, fmt.Errorf("could not get quotes: %w", err)
	}

	return quotes, nil
}

func (s *QuotesService) GetRandomQuote() (*models.Quote, error) {
	quotes, err := s.storage.GetQuotes(nil)
	if err != nil {
		return nil, fmt.Errorf("could not get quotes: %w", err)
	}

	if len(quotes) == 0 {
		return nil, errs.ErrQuoteNotFound
	}

	randomIndex := rand.Intn(len(quotes))

	return quotes[randomIndex], nil
}

func (s *QuotesService) DeleteQuote(id uint64) error {
	if id == 0 {
		return errs.ErrInvalidQuoteId
	}

	err := s.storage.DeleteQuote(id)
	if err != nil {
		return err
	}

	return nil
}
