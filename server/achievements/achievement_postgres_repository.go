package achievements

import (
	"context"
	"fmt"

	"github.com/abhinash-kml/nova/server/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostgresRepository struct {
	pgx    *pgxpool.Pool
	config *config.Config
	logger *zap.Logger
}

func NewPostgresRepository(pgx *pgxpool.Pool, c *config.Config, logger *zap.Logger) *PostgresRepository {
	return &PostgresRepository{
		pgx:    pgx,
		config: c,
		logger: logger,
	}
}

func (r *PostgresRepository) Create(ctx context.Context, req CreateAchievementRequest) (AchievementResponse, error) {
	query := `INSERT INTO
				achievements (
					key,
					name,
					description
				)
				VALUES (
					$1,
					$2,
					$3
				)
				RETURNING
					id,
					key,
					name,
					description`

	var response AchievementResponse

	err := r.pgx.QueryRow(
		ctx,
		query,
		req.Key,
		req.Name,
		req.Description,
	).Scan(
		&response.Id,
		&response.Key,
		&response.Name,
		&response.Description,
	)

	if err != nil {
		return AchievementResponse{}, fmt.Errorf("scanning row: %w", err)
	}

	return response, nil
}

func (r *PostgresRepository) Get(ctx context.Context, req GetAchievementRequest) (AchievementResponse, error) {
	query := `SELECT
				id,
				key,
				name,
				description
			FROM
				achievements
			WHERE
				id = $1`

	var response AchievementResponse

	err := r.pgx.QueryRow(
		ctx,
		query,
		req.Id,
	).Scan(
		&response.Id,
		&response.Key,
		&response.Name,
		&response.Description,
	)

	if err != nil {
		return AchievementResponse{}, fmt.Errorf("scanning row: %w", err)
	}

	return response, nil
}

func (r *PostgresRepository) Update(ctx context.Context, req UpdateAchievementRequest) error {
	query := `UPDATE
				achievements
			SET
				key = $1,
				name = $2,
				description = $3
			WHERE
				id = $4`

	_, err := r.pgx.Exec(
		ctx,
		query,
		req.Key,
		req.Name,
		req.Description,
		req.Id,
	)

	if err != nil {
		return fmt.Errorf("updating row: %w", err)
	}

	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, req DeleteAchievementRequest) error {
	query := `DELETE FROM
				achievements
			WHERE
				id = $1`

	_, err := r.pgx.Exec(
		ctx,
		query,
		req.Id,
	)
	if err != nil {
		return fmt.Errorf("delete row: %w", err)
	}

	return nil
}
