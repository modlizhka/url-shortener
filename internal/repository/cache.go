package repository

import (
	"sync"

	"url-shortener/pkg/storage"
)

type CacheStorage struct {
	data map[string]string
	sync.Mutex
}

func NewCacheStorage() *CacheStorage {
	return &CacheStorage{data: make(map[string]string)}
}

func (c *CacheStorage) GetLongUrl(shortURL string) (string, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	res, ok := c.data[shortURL]
	if !ok {
		return "", storage.ErrNotFound
	}

	return res, nil
}

func (s *CacheStorage) Insert(shortURL, longURL string) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	if _, ok := s.data[shortURL]; ok {
		return storage.ErrAlreadyExists
	}
	s.data[shortURL] = longURL
	return nil
}
