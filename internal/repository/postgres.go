package repository

import (
	"context"

	"url-shortener/pkg/storage"
	"url-shortener/pkg/storage/postgres"
)

type DataBaseStorage struct {
	pool *postgres.Pool
}

func NewDataBaseStorage(pool *postgres.Pool) *DataBaseStorage {
	return &DataBaseStorage{pool: pool}
}

func (s *DataBaseStorage) Insert(shortURL, longURL string) error {
	query := "INSERT INTO urls (short_url, long_url) VALUES ($1, $2) ON CONFLICT (short_url) DO NOTHING"
	_, err := s.pool.Query(context.Background(), query, shortURL, longURL)
	if postgres.IsDuplicateError(err) {
		return storage.ErrAlreadyExists
	}
	return err
}

func (s *DataBaseStorage) GetLongUrl(shortURL string) (string, error) {
	var longURL string
	err := s.pool.QueryRow(context.Background(), "SELECT long_url FROM urls WHERE short_url = $1", shortURL).Scan(&longURL)
	if err == postgres.ErrNotFound {
		return "", storage.ErrNotFound
	}
	return longURL, err
}
