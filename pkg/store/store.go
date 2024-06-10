package store

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
)

type Store struct {
	DB *pgx.Conn
}

// NewStorage creates a postgresql connection with pgx
func NewStore() (*Store, error) {
	connStr := viper.Get("postgres").(string)
	conn, err := pgx.Connect(context.Background(), connStr)

	if err != nil {
		return nil, err
	}

	return &Store{
		DB: conn,
	}, nil
}
