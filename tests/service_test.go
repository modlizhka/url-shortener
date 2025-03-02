package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"url-shortener/internal/service"
	"url-shortener/pkg/storage"

	"url-shortener/tests/mocks"

	"github.com/stretchr/testify/mock"
)

func TestShortening(t *testing.T) {
	// Создаём мок для Storage
	mockStorage := new(mocks.MockStorage)

	// Настройка ожидаемого поведения мока
	mockStorage.On("GetLongUrl", mock.Anything, mock.Anything).Return("", storage.ErrNotFound).Once()
	mockStorage.On("Insert", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	// Инициализируем сервис с использованием мока
	service := service.NewShortenerService(mockStorage)

	// Вызываем метод Shortening
	shortURL, err := service.Shortening(context.Background(), "https://example.com")
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
	mockStorage.On("GetLongUrl", context.Background(), "test_short_url").Return("https://example.com", nil).Once()

	service := service.NewShortenerService(mockStorage)

	longURL, err := service.Expansion(context.Background(), "test_short_url")
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", longURL)

	mockStorage.AssertExpectations(t)
}
