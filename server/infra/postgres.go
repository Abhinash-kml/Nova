package infra

import (
	"context"
	"database/sql"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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

func NewPostgresPGX(ctx context.Context, dsn string) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		zap.L().Error("Failed to connect to postgress using pgx library", zap.Error(err))
	}

	return conn
}

func NewPostgressPgxPool(ctx context.Context, dsn string) *pgxpool.Pool {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		zap.L().Error("Failed to parse pgx poll config", zap.Error(err))
		return nil
	}

	cfg.ConnConfig.Tracer = otelpgx.NewTracer()

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		zap.L().Error("Failed to connect to postgress pool", zap.Error(err))
		return nil
	}

	return pool
}
