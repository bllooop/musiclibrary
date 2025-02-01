package repository

import (
	"context"
	"fmt"

	logger "github.com/bllooop/musiclibrary/pkg/logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBname   string
	SSLMode  string
}

const (
	userListTable  = "userlist"
	songsListTable = "songlist"
)

func NewPostgresDB(cfg Config) (*pgxpool.Pool, error) {
	logger.Log.Info().Msg("Connecting to the database")
	db, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBname, cfg.SSLMode))
	if err != nil {
		return nil, err
	}
	return db, nil
}
