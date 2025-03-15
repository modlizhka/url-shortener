package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"url-shortener/internal/repository"
	"url-shortener/pkg/storage"
)

func TestCache_InsertAndGet(t *testing.T) {
	cache := repository.NewCacheStorage()

	// Тестируем вставку значения
	err := cache.Insert("test_key", "test_value")
	assert.NoError(t, err)

	// Тестируем получение существующего значения
	value, err := cache.GetLongUrl("test_key")
	assert.NoError(t, err)
	assert.Equal(t, "test_value", value)

	// Тестируем получение несуществующего значения
	_, err = cache.GetLongUrl("nonexistent_key")
	assert.Equal(t, storage.ErrNotFound, err)
}
