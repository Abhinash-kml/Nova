package achievements

import (
	"context"
	"fmt"
	"time"

	"github.com/abhinash-kml/nova/server/config"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type ProgressPostgresRepository struct {
	pgx    *pgxpool.Pool
	config *config.Config
	logger *zap.Logger
}

func NewProgressPostgresRepository(pgx *pgxpool.Pool, c *config.Config, logger *zap.Logger) *ProgressPostgresRepository {
	return &ProgressPostgresRepository{
		pgx:    pgx,
		config: c,
		logger: logger,
	}
}

func (r *ProgressPostgresRepository) Create(ctx context.Context, dto CreateProgressDTO) (AchievementProgress, error) {
	var result AchievementProgress

	query := `INSERT INTO 
			achievement_progress 
			(
				user_id, 
				criterion_id, 
				current_value, 
				updated_at
			)
			VALUES ($1, $2, $3, $4)
			RETURNING 
				id, 
				user_id, 
				criterion_id, 
				current_value, 
				updated_at;`

	now := time.Now()
	row := r.pgx.QueryRow(ctx, query, dto.UserId, dto.CriterionId, dto.CurrentValue, now)

	err := row.Scan(
		&result.Id,
		&result.UserId,
		&result.CriterionId,
		&result.CurrentValue,
		&result.UpdatedAt,
	)
	if err != nil {
		return AchievementProgress{}, fmt.Errorf("scanning row: %w", err)
	}

	return result, nil
}

func (r *ProgressPostgresRepository) Get(ctx context.Context, dto GetProgressDTO) (AchievementProgress, error) {
	var result AchievementProgress

	query := `SELECT 
				id, 
				user_id, 
				criterion_id, 
				current_value, 
				updated_at
		FROM 
			achievement_progress
		WHERE id = $1`

	row := r.pgx.QueryRow(ctx, query, dto.Id)
	err := row.Scan(
		&result.Id,
		&result.UserId,
		&result.CriterionId,
		&result.CurrentValue,
		&result.UpdatedAt,
	)
	if err != nil {
		return AchievementProgress{}, fmt.Errorf("scanning row: %w", err)
	}

	return result, nil
}

func (r *ProgressPostgresRepository) GetOfUser(ctx context.Context, dto GetProgressOfUserDTO) ([]AchievementProgress, error) {
	userId, err := uuid.Parse(dto.UserId)
	if err != nil {
		return nil, fmt.Errorf("parsing userid: %w", err)
	}

	query := `SELECT 
				id, 
				user_id, 
				criterion_id, 
				current_value, 
				updated_at
		FROM 
			achievement_progress
		WHERE user_id = $1
		ORDER BY id;`

	rows, err := r.pgx.Query(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("getting progress: %w", err)
	}
	defer rows.Close()

	result := make([]AchievementProgress, 0)

	for rows.Next() {
		var progress AchievementProgress

		err := rows.Scan(
			&progress.Id,
			&progress.UserId,
			&progress.CriterionId,
			&progress.CurrentValue,
			&progress.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		result = append(result, progress)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return result, nil
}

func (r *ProgressPostgresRepository) Update(ctx context.Context, dto UpdateProgressDTO) error {
	query := `UPDATE 
				achievement_progress
			SET 
				current_value = $1, 
				updated_at = $2
			WHERE id = $3`

	now := time.Now()
	_, err := r.pgx.Exec(ctx, query, dto.CurrentValue, now, dto.Id)
	if err != nil {
		return fmt.Errorf("updating progress: %w", err)
	}

	return nil
}

func (r *ProgressPostgresRepository) Delete(ctx context.Context, dto DeleteProgressDTO) (ProgressId, error) {
	query := `DELETE FROM 
				achievement_progress
			WHERE id = $1
			RETURNING id;`

	var deletedID int
	err := r.pgx.QueryRow(ctx, query, dto.Id).Scan(&deletedID)
	if err != nil {
		return dto.ProgressId, fmt.Errorf("scanning row: %w", err)
	}

	return ProgressId{Id: deletedID}, nil
}
