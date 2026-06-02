package infra

import (
	"database/sql"

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
