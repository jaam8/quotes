package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"quotes/internal/models"
	errs "quotes/pkg/errors"
	"testing"
)

type mockStorageAdapter struct {
	mock.Mock
}

func (m *mockStorageAdapter) SaveQuote(quote *models.Quote) (uint64, error) {
	args := m.Called(quote)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *mockStorageAdapter) GetQuotes(author *string) ([]*models.Quote, error) {
	args := m.Called(author)
	return args.Get(0).([]*models.Quote), args.Error(1)
}

func (m *mockStorageAdapter) GetQuote(id uint64) (*models.Quote, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Quote), args.Error(1)
}

func (m *mockStorageAdapter) DeleteQuote(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateQuote(t *testing.T) {
	testcases := []struct {
		name        string
		quote       *models.Quote
		exceptedId  uint64
		mockSetup   func(*mockStorageAdapter)
		exceptedErr error
	}{
		{
			name:       "success",
			exceptedId: 1,
			quote: &models.Quote{
				Author: "author",
				Quote:  "quote",
			},
			mockSetup: func(adapter *mockStorageAdapter) {
				quote := &models.Quote{
					Author: "author",
					Quote:  "quote",
				}
				adapter.On("SaveQuote", quote).Return(uint64(1), nil)
			},
			exceptedErr: nil,
		},
		{
			name: "invalid_quote",
			quote: &models.Quote{
				Author: "author",
				Quote:  "",
			},
			mockSetup:   func(adapter *mockStorageAdapter) {},
			exceptedErr: errs.ErrInvalidQuoteField,
		},
		{
			name: "invalid_author",
			quote: &models.Quote{
				Author: "",
				Quote:  "quote",
			},
			mockSetup:   func(adapter *mockStorageAdapter) {},
			exceptedErr: errs.ErrInvalidAuthorField,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			storageAdapter := new(mockStorageAdapter)
			tc.mockSetup(storageAdapter)

			service, err := NewQuotesService(storageAdapter)
			require.NoError(t, err)

			quoteId, err := service.CreateQuote(tc.quote)

			assert.ErrorIs(t, err, tc.exceptedErr)

			if tc.exceptedErr == nil {
				assert.NotNil(t, quoteId)
				assert.Equal(t, quoteId, tc.exceptedId)
			}

			storageAdapter.AssertExpectations(t)
		})

	}
}

func TestGetQuotes(t *testing.T) {
	testcases := []struct {
		name        string
		author      *string
		mockSetup   func(*mockStorageAdapter)
		exceptedErr error
	}{
		{
			name:   "success_by_author",
			author: func(s string) *string { return &s }("author"),
			mockSetup: func(adapter *mockStorageAdapter) {
				var author = "author"
				adapter.On("GetQuotes", &author).Return([]*models.Quote{}, nil)
			},
			exceptedErr: nil,
		},
		{
			name:   "success_by_id",
			author: nil,
			mockSetup: func(adapter *mockStorageAdapter) {
				var author *string
				adapter.On("GetQuotes", author).Return([]*models.Quote{}, nil)
			},
			exceptedErr: nil,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			storageAdapter := new(mockStorageAdapter)
			tc.mockSetup(storageAdapter)

			service, err := NewQuotesService(storageAdapter)
			require.NoError(t, err)

			quotes, err := service.GetQuotes(tc.author)
			assert.ErrorIs(t, err, tc.exceptedErr)

			if tc.exceptedErr == nil {
				assert.NotNil(t, quotes)
			}

			storageAdapter.AssertExpectations(t)
		})
	}
}

func TestGetRandomQuote(t *testing.T) {
	testcases := []struct {
		name        string
		mockSetup   func(*mockStorageAdapter)
		exceptedErr error
	}{
		{
			name: "success",
			mockSetup: func(adapter *mockStorageAdapter) {
				quotes := []*models.Quote{
					{Id: uint64(1), Author: "something author", Quote: "something quote"},
					{Id: uint64(2), Author: "something author", Quote: "something quote"},
					{Id: uint64(3), Author: "something author", Quote: "something quote"},
				}
				var author *string
				adapter.On("GetQuotes", author).Return(quotes, nil)
			},
			exceptedErr: nil,
		},
		{
			name: "quote_not_found",
			mockSetup: func(adapter *mockStorageAdapter) {
				var author *string
				adapter.On("GetQuotes", author).Return([]*models.Quote{}, nil)
			},
			exceptedErr: errs.ErrQuoteNotFound,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			storageAdapter := new(mockStorageAdapter)
			tc.mockSetup(storageAdapter)

			service, err := NewQuotesService(storageAdapter)
			require.NoError(t, err)

			quote, err := service.GetRandomQuote()
			assert.ErrorIs(t, err, tc.exceptedErr)

			if tc.exceptedErr == nil {
				anotherQuote, err := service.GetRandomQuote()
				assert.ErrorIs(t, err, tc.exceptedErr)

				for quote.Id == anotherQuote.Id {
					anotherQuote, err = service.GetRandomQuote()
					assert.ErrorIs(t, err, tc.exceptedErr)
				}
				assert.NotEqual(t, quote.Id, anotherQuote.Id)
			}

			storageAdapter.AssertExpectations(t)
		})
	}
}

func TestDeleteQuote(t *testing.T) {
	testcases := []struct {
		name        string
		id          uint64
		mockSetup   func(*mockStorageAdapter)
		exceptedErr error
	}{
		{
			name: "success",
			id:   1,
			mockSetup: func(adapter *mockStorageAdapter) {
				adapter.On("DeleteQuote", uint64(1)).Return(nil)
			},
			exceptedErr: nil,
		},
		{
			name:        "invalid_id",
			id:          0,
			mockSetup:   func(adapter *mockStorageAdapter) {},
			exceptedErr: errs.ErrInvalidQuoteId,
		},
		{
			name: "quote_not_found",
			id:   999,
			mockSetup: func(adapter *mockStorageAdapter) {
				adapter.On("DeleteQuote", uint64(999)).Return(errs.ErrQuoteNotFound)
			},
			exceptedErr: errs.ErrQuoteNotFound,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			storageAdapter := new(mockStorageAdapter)

			tc.mockSetup(storageAdapter)
			service, err := NewQuotesService(storageAdapter)
			require.NoError(t, err)

			err = service.DeleteQuote(tc.id)
			assert.ErrorIs(t, err, tc.exceptedErr)

			storageAdapter.AssertExpectations(t)
		})
	}
}
