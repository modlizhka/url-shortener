package mocks

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

// MockPool is a mock implementation of pgxpool.Pool
type MockPool struct {
	mock.Mock
}

// Query provides a mock function for pgxpool.Pool.Query
func (m *MockPool) Query(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	argsList := m.Called(ctx, query, args)
	return argsList.Get(0).(pgconn.CommandTag), argsList.Error(1)
}

// QueryRow provides a mock function for pgxpool.Pool.QueryRow
func (m *MockPool) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return m.Called(ctx, query, args).Get(0).(pgx.Row)
}
