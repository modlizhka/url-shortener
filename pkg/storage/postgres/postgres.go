package postgres

import (
	"context"
	"fmt"
	"time"

	"errors"

	"url-shortener/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	*pgxpool.Pool
}

var ErrNotFound = pgx.ErrNoRows

func NewClient(ctx context.Context, cfg config.DataBase) (Pool, error) {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	p, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return Pool{}, err
	}
	return Pool{p}, nil
}

func IsDuplicateError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return true
	}
	return false
}
