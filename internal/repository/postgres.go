package repository

import (
	"context"
	// "errors"
	// "strconv"

	"url-shortener/pkg/storage"
	"url-shortener/pkg/storage/postgres"
)

type DataBaseStorage struct {
	pool *postgres.Pool
}

func NewDataBaseStorage(pool *postgres.Pool) *DataBaseStorage {
	return &DataBaseStorage{pool: pool}
}

func (s *DataBaseStorage) Insert(ctx context.Context, shortURL, longURL string) error {
	query := "INSERT INTO urls (short_url, long_url) VALUES ($1, $2) ON CONFLICT (short_url) DO NOTHING"
	_, err := s.pool.Query(ctx, query, shortURL, longURL)
	if postgres.IsDuplicateError(err) {
		return storage.ErrAlreadyExists
	}
	return err
}

func (s *DataBaseStorage) GetLongUrl(ctx context.Context, shortURL string) (string, error) {
	var longURL string
	err := s.pool.QueryRow(ctx, "SELECT long_url FROM urls WHERE short_url = $1", shortURL).Scan(&longURL)
	if err == postgres.ErrNotFound {
		return "", storage.ErrNotFound
	}
	return longURL, err
}
