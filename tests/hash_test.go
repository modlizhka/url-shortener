package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"url-shortener/internal/service"
)

func TestEncodeHash(t *testing.T) {
	input := "https://example.com"
	result := service.EncodeHash(input)

	assert.Len(t, result, 8) // Проверяем длину результата
	for _, char := range result {
		assert.Contains(t, service.Alphabet, string(char)) // Проверяем допустимые символы
	}
}
