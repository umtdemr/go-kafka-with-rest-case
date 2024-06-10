package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
)

type Storage struct {
	DB *pgx.Conn
}

// NewStorage creates a postgresql connection with pgx
func NewStorage() (*Storage, error) {
	connStr := viper.Get("postgres").(string)
	conn, err := pgx.Connect(context.Background(), connStr)

	if err != nil {
		return nil, err
	}

	return &Storage{
		DB: conn,
	}, nil
}
