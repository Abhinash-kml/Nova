package infra

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func NewPostgres(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		zap.L().Fatal("Failed to connect to postgress", zap.Error(err))
	}

	return db
}

func NewPostgresPGX(dsn string) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		zap.L().Error("Failed to connect to postgress using pgx library", zap.Error(err))
	}

	return conn
}
