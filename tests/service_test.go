package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"url-shortener/internal/service"
	"url-shortener/pkg/storage"

	"url-shortener/tests/mocks"

	"github.com/stretchr/testify/mock"
)

func TestShortening(t *testing.T) {
	mockStorage := new(mocks.MockStorage)

	mockStorage.On("GetLongUrl", mock.Anything, mock.Anything).Return("", storage.ErrNotFound).Once()
	mockStorage.On("Insert", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	service := service.NewShortenerService(mockStorage)

	shortURL, err := service.Shortening("https://example.com")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Проверяем длину сгенерированного короткого URL
	if len(shortURL) != 10 {
		t.Errorf("expected length of 10, got %d", len(shortURL))
	}

	// Проверяем, что все ожидания были выполнены
	mockStorage.AssertExpectations(t)
}

func TestExpansion(t *testing.T) {
	mockStorage := new(mocks.MockStorage)
	mockStorage.On("GetLongUrl", "test_short_url").Return("https://example.com", nil).Once()

	service := service.NewShortenerService(mockStorage)

	longURL, err := service.Expansion("test_short_url")
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", longURL)

	mockStorage.AssertExpectations(t)
}
