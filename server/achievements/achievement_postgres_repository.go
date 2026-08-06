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

func (r *PostgresRepository) Create(ctx context.Context, dto CreateAchievementDTO) (AchievementResponseDTO, error) {
	var response AchievementResponseDTO

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

	err := r.pgx.QueryRow(
		ctx,
		query,
		dto.Key,
		dto.Name,
		dto.Description,
	).Scan(
		&response.Id,
		&response.Key,
		&response.Name,
		&response.Description,
	)

	if err != nil {
		return AchievementResponseDTO{}, fmt.Errorf("scanning row: %w", err)
	}

	return response, nil
}

func (r *PostgresRepository) Get(ctx context.Context, dto GetAchievementDTO) (AchievementResponseDTO, error) {
	var response AchievementResponseDTO

	query := `SELECT
				id,
				key,
				name,
				description
			FROM
				achievements
			WHERE
				id = $1`

	err := r.pgx.QueryRow(
		ctx,
		query,
		dto.Id,
	).Scan(
		&response.Id,
		&response.Key,
		&response.Name,
		&response.Description,
	)

	if err != nil {
		return AchievementResponseDTO{}, fmt.Errorf("scanning row: %w", err)
	}

	return response, nil
}

func (r *PostgresRepository) Update(ctx context.Context, dto UpdateAchievementDTO) error {
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
		dto.Key,
		dto.Name,
		dto.Description,
		dto.Id,
	)

	if err != nil {
		return fmt.Errorf("updating row: %w", err)
	}

	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, dto DeleteAchievementDTO) error {
	query := `DELETE FROM
				achievements
			WHERE
				id = $1`

	_, err := r.pgx.Exec(
		ctx,
		query,
		dto.Id,
	)

	if err != nil {
		return fmt.Errorf("delete row: %w", err)
	}

	return nil
}
