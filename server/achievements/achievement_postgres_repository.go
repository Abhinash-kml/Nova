package achievements

import (
	"context"
	"fmt"
	"uuid"

	"github.com/Masterminds/squirrel"
	"github.com/abhinash-kml/nova/server/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostgresRepository struct {
	pgx          *pgxpool.Pool
	config       *config.Config
	logger       *zap.Logger
	queryBuilder squirrel.StatementBuilderType
}

func NewPostgresRepository(pgx *pgxpool.Pool, c *config.Config, logger *zap.Logger) *PostgresRepository {
	return &PostgresRepository{
		pgx:          pgx,
		config:       c,
		logger:       logger,
		queryBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *PostgresRepository) Create(ctx context.Context, req CreateAchievementRequest) (AchievementResponse, error) {
	query := `INSERT INTO
				achievements (
					key,
					name,
					description,
					stat_id,
					target_value
				)
				VALUES (
					$1,
					$2,
					$3,
					$4,
					$5
				)
				RETURNING
					*;`

	var response AchievementResponse

	err := r.pgx.QueryRow(
		ctx,
		query,
		req.Key,
		req.Name,
		req.Description,
		req.Stat,
		req.TargetValue,
	).Scan(
		&response.Id,
		&response.Key,
		&response.Name,
		&response.Description,
		&response.Stat,
		&response.TargetValue,
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
				description,
				stat_id,
				target_value
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
		&response.Stat,
		&response.TargetValue,
	)

	if err != nil {
		return AchievementResponse{}, fmt.Errorf("scanning row: %w", err)
	}

	return response, nil
}

func (r *PostgresRepository) Update(ctx context.Context, req UpdateAchievementRequest) (AchievementResponse, error) {
	builder := r.queryBuilder.Update("achievements").Suffix("RETURNING *")
	if req.Name != nil {
		builder = builder.Set("name", req.Name)
	}
	if req.Key != nil {
		builder = builder.Set("key", req.Key)
	}
	if req.Description != nil {
		builder = builder.Set("description", req.Description)
	}
	if req.Stat != nil {
		builder = builder.Set("stat_id", req.Stat)
	}
	if req.TargetValue != nil {
		builder = builder.Set("target_value", req.TargetValue)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return AchievementResponse{}, fmt.Errorf("generating query: %w", err)
	}

	var a AchievementResponse
	err = r.pgx.QueryRow(ctx, query, args...).Scan(
		&a.Id,
		&a.Key,
		&a.Name,
		&a.Description,
		&a.Stat,
		&a.TargetValue,
	)
	if err != nil {
		return AchievementResponse{}, fmt.Errorf("scanning row: %w", err)
	}

	return a, nil
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

func (r *PostgresRepository) CompleteAchievement(ctx context.Context, req CompleteAchievementRequest) error {
	rawQuery := `INSERT INTO 
					completed_achievements(
						user_id,
						achievement_id,
						completed_on
					) VALUES (
					 	$1,
						$2,
						CURRENT_TIMESTAMP
					);`

	_, err := r.pgx.Exec(ctx, rawQuery, req.UserId, req.AchievementId)
	if err != nil {
		return fmt.Errorf("create completed achievement: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetCompletedAchievementsOfUser(ctx context.Context, req CompletedAchievementRequest) ([]int, error) {
	rawQuery := `SELECT
					achievement_id
				FROM
					completed_achievements
				WHERE
					user_id = $1;	
	`

	rows, err := r.pgx.Query(ctx, rawQuery, req.UserId)
	if err != nil {
		return nil, fmt.Errorf("getting achievements: %w", err)
	}
	defer rows.Close()

	var achievements []int
	var achievement int
	for rows.Next() {
		if err = rows.Scan(&achievement); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}
		achievements = append(achievements, achievement)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return achievements, nil
}

func (r *PostgresRepository) GetProgress(ctx context.Context, dto GetProgressRequest) (AchievementProgress, error) {
	query := `SELECT 
				achievement_id, 
				current_stat_value, 
				updated_at
			FROM 
				achievement_progress
			WHERE 
				id = $1;`

	var result AchievementProgress

	row := r.pgx.QueryRow(ctx, query, dto.Id)
	err := row.Scan(
		&result.AchievementId,
		&result.CurrentValue,
		&result.UpdatedAt,
	)
	if err != nil {
		return AchievementProgress{}, fmt.Errorf("scanning row: %w", err)
	}

	return result, nil
}

func (r *PostgresRepository) GetProgressOfUser(ctx context.Context, dto GetProgressOfUserRequest) ([]AchievementProgress, error) {
	userId, err := uuid.Parse(dto.UserId)
	if err != nil {
		return nil, fmt.Errorf("parsing userid: %w", err)
	}

	query := `SELECT 
				achievement_id,
				current_stat_value, 
				updated_at
			FROM 
				achievement_progress
			WHERE 
				user_id = $1
			ORDER BY 
				achievement_id;`

	rows, err := r.pgx.Query(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("getting progress: %w", err)
	}
	defer rows.Close()

	var result []AchievementProgress
	var progress AchievementProgress

	for rows.Next() {
		err := rows.Scan(
			&progress.AchievementId,
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

func (r *PostgresRepository) UpdateProgress(ctx context.Context, dto UpdateProgressRequest) error {
	query := `INSERT INTO
				achievement_progress(
					achievement_id,
					user_id,
					current_stat_value,
					updated_at
				) VALUES (
					$1,
					$2,
					$3,
					CURRENT_TIMESTAMP
				)
			ON ONFLICT (
				achievement_id,
				user_id
			)
			DO UPDATE SET
				current_stat_value = EXCLUDED.stat_current_value,
				updated_at = CURRENT_TIMESTAMP;
	`

	_, err := r.pgx.Exec(ctx, query, dto.AchievementId, dto.UserId, dto.CurrentValue)
	if err != nil {
		return fmt.Errorf("updating progress: %w", err)
	}

	return nil
}

func (r *PostgresRepository) DeleteProgress(ctx context.Context, dto DeleteProgressRequest) (ProgressId, error) {
	query := `DELETE FROM 
				achievement_progress
			WHERE 
				id = $1
			RETURNING 
				id;`

	var deletedID int

	err := r.pgx.QueryRow(ctx, query, dto.Id).Scan(&deletedID)
	if err != nil {
		return dto.ProgressId, fmt.Errorf("scanning row: %w", err)
	}

	return ProgressId{Id: deletedID}, nil
}

func (r *PostgresRepository) DeleteAllProgressOfUser(ctx context.Context, req DeleteProgressOfUserRequest) error {
	query := `
		DELETE FROM
			achievement_progress
		WHERE
			user_id = $1;
	`

	_, err := r.pgx.Exec(ctx, query, req.UserId)
	if err != nil {
		return fmt.Errorf("deleting user progrogress: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetAchievementsForEvaluator(ctx context.Context) ([]Achievement, error) {
	query := `
		SELECT
			id,
			stat_id,
			target_value
		FROM
			achievements;
	`

	rows, err := r.pgx.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("getting achievements: %w", err)
	}
	defer rows.Close()

	var achievements []Achievement
	var achievement Achievement

	for rows.Next() {
		if err = rows.Scan(&achievement.AchievementId, &achievement.StatId, &achievement.TargetValue); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}
		achievements = append(achievements, achievement)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return achievements, nil
}
