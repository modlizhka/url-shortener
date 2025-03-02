package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"url-shortener/pkg/storage"
)

const Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_"
const hashLength = 8 // Длина желаемого хэша
const maxIndex = len(Alphabet) ^ 2 - 1

//go:generate mockgen -source=user_service.go -destination=mocks/mock.go
type Storage interface {
	GetLongUrl(ctx context.Context, shortUrl string) (string, error)
	Insert(ctx context.Context, shortUrl, longUrl string) error
}

type ShortenerService struct {
	Storage Storage
}

func NewShortenerService(Storage Storage) *ShortenerService {
	return &ShortenerService{Storage: Storage}
}

func (s ShortenerService) Shortening(ctx context.Context, longUrl string) (string, error) {
	id := 0
	hash := EncodeHash(longUrl)
	for id < maxIndex {
		shortUrl := hash + IntToIndex63(id)
		longCheck, err := s.Storage.GetLongUrl(ctx, shortUrl)
		if err == storage.ErrNotFound {
			err = s.Storage.Insert(ctx, shortUrl, longUrl)
			return shortUrl, err
		} else if longCheck == longUrl {
			return shortUrl, err
		}
		id++
	}
	return "", errors.New("index out of range")
}

func (s ShortenerService) Expansion(ctx context.Context, shortUrl string) (string, error) {
	res, err := s.Storage.GetLongUrl(ctx, shortUrl)
	return res, err
}

// Функция для преобразования байтов в строку фиксированной длины
func EncodeHash(input string) string {
	// Создаем хэш
	hasher := sha256.New()
	hasher.Write([]byte(input))
	hash := hasher.Sum(nil)

	result := ""
	for i := 0; i < hashLength; i++ {
		index := int(hash[i]) % len(Alphabet) // Получаем индекс из хэша
		result += string(Alphabet[index])     // Добавляем символ из алфавита
	}
	return result
}

// Функция для преобразования числа из 10-тичной системы счисления в 63-ную
func IntToIndex63(id int) string {
	res := ""
	res += string(Alphabet[id/len(Alphabet)])
	res += string(Alphabet[id%len(Alphabet)])
	return res

}
