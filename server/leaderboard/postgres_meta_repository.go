package leaderboard

import (
	"context"
	"time"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostgresMetaRepository struct {
	logger *zap.Logger
	pgx    *pgxpool.Pool
}

func NewPostgresMetaRepository(p *pgxpool.Pool, l *zap.Logger) *PostgresMetaRepository {
	return &PostgresMetaRepository{
		pgx:    p,
		logger: l,
	}
}

func (r *PostgresMetaRepository) GetAll(ctx context.Context, cursor int, limit int) ([]Leaderboard, error) {
	var rows pgx.Rows
	var err error

	rawQuery := `SELECT
					id,
					name,
					type,
					process_interval,
					created_by,
					created_at
				FROM
					leaderboards`

	if cursor == 0 {
		rawQuery += ` LIMIT
						$1;`
		rows, err = r.pgx.Query(ctx, rawQuery, limit)
	} else {
		rawQuery += ` WHERE
						id > $1
					LIMIT
						$2;`

		rows, err = r.pgx.Query(ctx, rawQuery, cursor, limit)
	}
	if err != nil {
		r.logger.Error("Failed to execute getall query", zap.Error(err))
		return nil, common.TranslatePostgresError(err, r.logger)
	}

	var leaderboards []Leaderboard

	// Scan returned rows
	for rows.Next() {
		var leaderboard Leaderboard
		err := rows.Scan(&leaderboard.Id, &leaderboard.Name, &leaderboard.Type, &leaderboard.ProcessInterval,
			&leaderboard.CreatedBy, &leaderboard.CreatedAt)
		if err != nil {
			r.logger.Error("failed to scan returned row in getall query", zap.Error(err))
			return nil, common.TranslatePostgresError(err, r.logger)
		}
	}

	return leaderboards, nil
}

func (r *PostgresMetaRepository) Get(ctx context.Context, id uuid.UUID) (Leaderboard, error) {
	rawQuery := `SELECT
					id,
					name,
					type,
					process_interval,
					created_by,
					created_at
				FROM
					leaderboards
				WHERE
					id = $1;`

	// Execute query
	row := r.pgx.QueryRow(ctx, rawQuery)

	var leaderboard Leaderboard

	// Scan returned row
	err := row.Scan(&leaderboard.Id, &leaderboard.Name, &leaderboard.Type, &leaderboard.ProcessInterval,
		&leaderboard.CreatedBy, &leaderboard.CreatedAt)
	if err != nil {
		r.logger.Error("Failed to scan returned row in getbyid query", zap.Error(err))
		return Leaderboard{}, common.TranslatePostgresError(err, r.logger)
	}

	return leaderboard, nil
}

func (r *PostgresMetaRepository) Create(ctx context.Context, dto CreateDTO) (Leaderboard, error) {
	rawQuery := `INSERT INTO
					leaderboards(id, name, type, process_interval, created_by, created_at)
				VALUES($1, $2, $3, $4, $5, $6);`

	id, _ := uuid.NewV7()
	now := time.Now()

	// Execute query
	_, err := r.pgx.Exec(ctx, rawQuery, id, dto.Name, dto.Type, dto.ProcessInterval, dto.CreatedBy, now)
	if err != nil {
		r.logger.Error("Failed to execute insert query", zap.Error(err))
		return Leaderboard{}, common.TranslatePostgresError(err, r.logger)
	}
	createdBy, _ := uuid.Parse(dto.CreatedBy)

	return Leaderboard{
		Id:              id,
		Name:            dto.Name,
		Type:            dto.Type,
		ProcessInterval: dto.ProcessInterval,
		CreatedBy:       createdBy,
		CreatedAt:       now,
	}, nil
}

func (r *PostgresMetaRepository) Modify(ctx context.Context, dto ModifyDTO) (Leaderboard, error) {
	rawQuery := `UPDATE
					leaderboards
				SET
					type = $2,
					process_interval = $3
				WHERE
					id = $1
				RETUNRING
					*;`

	// Execute query
	row := r.pgx.QueryRow(ctx, rawQuery, dto.Id, dto.Type, dto.ProcessInterval)

	var leaderboard Leaderboard

	// Scan returned row
	err := row.Scan(&leaderboard.Id, &leaderboard.Name, &leaderboard.Type, &leaderboard.ProcessInterval,
		&leaderboard.CreatedBy, &leaderboard.CreatedAt)
	if err != nil {
		r.logger.Error("Failed to scan returned row in modify query", zap.Error(err))
		return Leaderboard{}, common.TranslatePostgresError(err, r.logger)
	}

	return leaderboard, nil
}

func (r *PostgresMetaRepository) Delete(ctx context.Context, dto DeleteDTO) (Leaderboard, error) {
	rawQuery := `DELETE FROM
					leaderboards
				WHERE
					id = $1
				RETURNING
					*;`

	// Execute query
	row := r.pgx.QueryRow(ctx, rawQuery, dto.Id)

	var leaderboard Leaderboard

	// Scan returned row
	err := row.Scan(&leaderboard.Id, &leaderboard.Name, &leaderboard.Type, &leaderboard.ProcessInterval,
		&leaderboard.CreatedBy, &leaderboard.CreatedAt)
	if err != nil {
		r.logger.Error("Failed to scan returned row in modify query", zap.Error(err))
		return Leaderboard{}, common.TranslatePostgresError(err, r.logger)
	}

	return leaderboard, nil
}
