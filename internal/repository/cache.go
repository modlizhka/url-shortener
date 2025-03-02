package repository

import (
	"sync"
	"time"

	"url-shortener/pkg/storage"
)

type Node struct {
	data        string
	RequestTime time.Time
}

type CacheStorage struct {
	data map[string]Node
	sync.Mutex
}

func NewCacheStorage() *CacheStorage {
	return &CacheStorage{data: make(map[string]Node)}
}

func (c *CacheStorage) GetLongUrl(shortURL string) (string, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	res, ok := c.data[shortURL]
	if !ok {
		return "", storage.ErrNotFound
	}
	c.data[shortURL] = Node{
		data:        res.data,
		RequestTime: time.Now(),
	}
	return res.data, nil
}

func (s *CacheStorage) Insert(shortURL, longURL string) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	if _, ok := s.data[shortURL]; ok {
		return storage.ErrAlreadyExists
	}
	s.data[shortURL] = Node{
		data:        longURL,
		RequestTime: time.Now(),
	}
	return nil
}

func (s *CacheStorage) CashChecker(lifeTime int64) {
	for {
		time.Sleep(time.Duration(lifeTime) * time.Millisecond)
		s.Mutex.Lock()
		for key := range s.data {
			if int64(time.Now().Sub(s.data[key].RequestTime).Minutes()) > lifeTime {
				delete(s.data, key)
			}
		}
		s.Unlock()
	}
}
