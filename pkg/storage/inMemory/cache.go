package inmemory

import (
	"sync"
	"url-shortener/pkg/storage"
)

type MemoryStorage struct {
	data map[string]string
	sync.Mutex
}

func NewStorage(hasher func(string) string) *MemoryStorage {
	return &MemoryStorage{data: make(map[string]string)}
}

func (s *MemoryStorage) Insert(key, value string) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	_, ok := s.data[key]
	if !ok {
		return storage.ErrAlreadyExists
	}
	s.data[key] = value
	return nil
}

func (s *MemoryStorage) Get(key string) (string, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	res, ok := s.data[key]
	if !ok {
		return "", storage.ErrNotFound
	}
	return res, nil
}
