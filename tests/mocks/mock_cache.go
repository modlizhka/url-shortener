// Automatically generated by MockGen. DO NOT EDIT!
// Source: internal/repository/cache.go

package mocks

import (
    // "url-shortener/pkg/storage"

    "github.com/stretchr/testify/mock"
)

// MockCacheStorage is an autogenerated mock type for the CacheStorage type
type MockCacheStorage struct {
    mock.Mock
}

// GetLongUrl provides a mock function with given fields: shortURL
func (_m *MockCacheStorage) GetLongUrl(shortURL string) (string, error) {
    ret := _m.Called(shortURL)

    var r0 string
    if rf, ok := ret.Get(0).(func(string) string); ok {
        r0 = rf(shortURL)
    } else {
        r0 = ret.Get(0).(string)
    }

    var r1 error
    if rf, ok := ret.Get(1).(func(string) error); ok {
        r1 = rf(shortURL)
    } else {
        r1 = ret.Error(1)
    }

    return r0, r1
}

// Insert provides a mock function with given fields: shortURL, longURL
func (_m *MockCacheStorage) Insert(shortURL, longURL string) error {
    ret := _m.Called(shortURL, longURL)

    var r0 error
    if rf, ok := ret.Get(0).(func(string, string) error); ok {
        r0 = rf(shortURL, longURL)
    } else {
        r0 = ret.Error(0)
    }

    return r0
}
