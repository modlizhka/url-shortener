package mocks

import (

	// "url-shortener/pkg/storage"

	"github.com/stretchr/testify/mock"
)

type MockShortenerService struct {
	mock.Mock
}

func (m *MockShortenerService) Shortening(longUrl string) (string, error) {
	args := m.Called(longUrl)
	return args.String(0), args.Error(1)
}

func (m *MockShortenerService) Expansion(shortUrl string) (string, error) {
	args := m.Called(shortUrl)
	return args.String(0), args.Error(1)
}
