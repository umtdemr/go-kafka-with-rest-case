package store

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
)

type Store struct {
	db *pgx.Conn
}

type CreateLogData struct {
	Operation   string `json:"operation"`
	RequestTime int    `json:"request_time"`
	Timestamp   int    `json:"timestamp"`
}

// NewStore creates a postgresql connection with pgx
func NewStore() (*Store, error) {
	connStr := viper.Get("postgres").(string)
	conn, err := pgx.Connect(context.Background(), connStr)

	if err != nil {
		return nil, err
	}

	return &Store{
		db: conn,
	}, nil
}

func (s *Store) Init() error {
	query := `CREATE TABLE IF NOT EXISTS "api_logs" (
		id serial PRIMARY KEY,
		operation varchar(10) NOT NULL,
		request_time int NOT NULL,
		timestamp BIGINT NOT NULL
	)`
	_, err := s.db.Exec(context.Background(), query)

	return err
}

func (s *Store) Close() {
	s.db.Close(context.Background())
}

func (s *Store) CreateLog(data *CreateLogData) error {
	query := `INSERT INTO "api_logs"(operation, request_time, timestamp) VALUES (@operation, @requestTime, @timestamp)`
	args := pgx.NamedArgs{
		"operation":   data.Operation,
		"requestTime": data.RequestTime,
		"timestamp":   data.Timestamp,
	}

	_, err := s.db.Exec(context.Background(), query, args)
	return err
}

func (s *Store) GetAllLogs() ([]CreateLogData, error) {
	query := `SELECT operation, request_time, timestamp FROM "api_logs"`
	rows, err := s.db.Query(context.Background(), query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	allLogs := []CreateLogData{}

	for rows.Next() {
		var singleLog CreateLogData

		err := rows.Scan(&singleLog.Operation, &singleLog.RequestTime, &singleLog.Timestamp)

		if err != nil {
			return nil, err
		}
		allLogs = append(allLogs, singleLog)
	}

	return allLogs, nil
}
