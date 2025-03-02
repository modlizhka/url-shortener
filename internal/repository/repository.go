package repository

import (
	"context"
	"url-shortener/pkg/storage"
)

type Storage struct {
	db    *DataBaseStorage
	cache *CacheStorage
}

func NewStorage(db *DataBaseStorage, cache *CacheStorage) *Storage {
	return &Storage{db: db, cache: cache}
}

func (s *Storage) GetLongUrl(ctx context.Context, shortUrl string) (string, error) {
	res, err := s.cache.GetLongUrl(shortUrl)
	if err != nil {
		if err == storage.ErrNotFound {
			res, err = s.db.GetLongUrl(ctx, shortUrl)
			if err == nil {
				s.cache.Insert(shortUrl, res)
			}
		} else {
			res, err = s.db.GetLongUrl(ctx, shortUrl)
		}
	}
	return res, err

}

func (s *Storage) Insert(ctx context.Context, shortURL, longURL string) error {
	err := s.db.Insert(ctx, shortURL, longURL)
	if err != nil {
		return err
	}
	s.cache.Insert(shortURL, longURL)
	return nil
}
