package achievements

import (
	"context"
	"time"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type ProgressPostgresRepository struct {
	pgx    *pgxpool.Pool
	logger *zap.Logger
}

func NewProgressPostgresRepository(pgx *pgxpool.Pool, logger *zap.Logger) *ProgressPostgresRepository {
	return &ProgressPostgresRepository{
		pgx:    pgx,
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
		r.logger.Error("failed to scan create achievement progress query", zap.Error(err))
		return AchievementProgress{}, common.TranslatePostgresError(err, r.logger)
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
		r.logger.Error("failed to scan row in get achievement progress query", zap.Error(err))
		return AchievementProgress{}, common.TranslatePostgresError(err, r.logger)
	}

	return result, nil
}

func (r *ProgressPostgresRepository) GetOfUser(ctx context.Context, dto GetProgressOfUserDTO) ([]AchievementProgress, error) {
	userId, err := uuid.Parse(dto.UserId)
	if err != nil {
		r.logger.Error("failed to parse user id", zap.Error(err))
		return nil, common.TranslatePostgresError(err, r.logger)
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
		r.logger.Error("failed to execute get achievement of user query", zap.Error(err))
		return nil, common.TranslatePostgresError(err, r.logger)
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
			r.logger.Error("failed to scan achievement progress in ", zap.Error(err))
			return nil, common.TranslatePostgresError(err, r.logger)
		}

		result = append(result, progress)
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
		r.logger.Error("failed to update achievement progress", zap.Error(err))
		return common.TranslatePostgresError(err, r.logger)
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
		r.logger.Error("failed to execute delete achievement progress query", zap.Error(err))
		return dto.ProgressId, common.TranslatePostgresError(err, r.logger)
	}

	return ProgressId{Id: deletedID}, nil
}
