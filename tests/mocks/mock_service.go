package mocks

import (
	"context"

	// "url-shortener/pkg/storage"

	"github.com/stretchr/testify/mock"
)

type MockShortenerService struct {
	mock.Mock
}

func (m *MockShortenerService) Shortening(ctx context.Context, longUrl string) (string, error) {
	args := m.Called(ctx, longUrl)
	return args.String(0), args.Error(1)
}

func (m *MockShortenerService) Expansion(ctx context.Context, shortUrl string) (string, error) {
	args := m.Called(ctx, shortUrl)
	return args.String(0), args.Error(1)
}
