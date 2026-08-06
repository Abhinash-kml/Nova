package achievements

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/abhinash-kml/nova/server/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type CriteriaPostgresRepository struct {
	logger           *zap.Logger
	config           *config.Config
	pgx              *pgxpool.Pool
	statementBuilder squirrel.StatementBuilderType
}

func NewCriteriaPostgresRepository(connection *pgxpool.Pool, c *config.Config, l *zap.Logger) *CriteriaPostgresRepository {
	return &CriteriaPostgresRepository{
		logger:           l,
		config:           c,
		pgx:              connection,
		statementBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *CriteriaPostgresRepository) Add(ctx context.Context, dto CreateCriteriaDTO) (AchievementCriteria, error) {
	rawQuery := `INSERT INTO
					achievement_criteria(
						achievement_id,
						stat_id,
						target_value,
						criteria_text
					)
				VALUES ($1, $2, $3, $4)
				RETURNING
					achievement_id,
					stat_id,
					target_value,
					criteria_text;`

	var criteria AchievementCriteria

	err := r.pgx.QueryRow(
		ctx,
		rawQuery,
		dto.AchievementId,
		dto.StatId,
		dto.TargetValue,
		dto.CriteriaText,
	).Scan(
		&criteria.AchievementId,
		&criteria.StatId,
		&criteria.TargetValue,
		&criteria.CriteriaText,
	)

	if err != nil {
		return AchievementCriteria{}, fmt.Errorf("scanning row: %w", err)
	}

	return criteria, nil
}

func (r *CriteriaPostgresRepository) GetAll(ctx context.Context, dto GetAllCriteriaDTO) ([]AchievementCriteria, error) {
	var rows pgx.Rows
	var err error

	rawQuery := `SELECT 
					achievement_id,
					stat_id,
					target_value,
					criteria_text
				FROM
					achievement_criteria`
	if dto.Cursor == 0 {
		rawQuery += ` LIMIT $1;`

		rows, err = r.pgx.Query(ctx, rawQuery, dto.Limit)
	} else {
		rawQuery += ` WHERE id > $1
					AND LIMIT $2;`

		rows, err = r.pgx.Query(ctx, rawQuery, dto.Cursor, dto.Limit)
	}
	if err != nil {
		return []AchievementCriteria{}, fmt.Errorf("get all criteria: %w", err)
	}

	defer rows.Close()

	var criterias []AchievementCriteria

	for rows.Next() {
		var criteria AchievementCriteria
		err = rows.Scan(
			&criteria.AchievementId,
			&criteria.StatId,
			&criteria.TargetValue,
			&criteria.CriteriaText,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		criterias = append(criterias, criteria)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return criterias, nil
}

func (r *CriteriaPostgresRepository) Get(ctx context.Context, dto GetCriteriaDTO) (AchievementCriteria, error) {
	rawQuery := `SELECT
					achievement_id,
					stat_id,
					target_value,
					criteria_text
				FROM
					achievement_criteria
				WHERE
					id = $1`

	var criteria AchievementCriteria

	err := r.pgx.QueryRow(
		ctx,
		rawQuery,
		dto.CriteriaID,
	).Scan(
		&criteria.AchievementId,
		&criteria.StatId,
		&criteria.TargetValue,
		&criteria.CriteriaText,
	)

	if err != nil {
		return AchievementCriteria{}, fmt.Errorf("scanning row: %w", err)
	}

	return criteria, nil
}

func (r *CriteriaPostgresRepository) GetByAchievement(ctx context.Context, dto GetCriteriaByAchievementDTO) ([]AchievementCriteria, error) {
	rawQuery := `SELECT
					achievement_id,
					stat_id,
					target_value,
					criteria_text
				FROM
					achievement_criteria
				WHERE
					achievement_id = $1
				ORDER BY
					stat_id;`

	rows, err := r.pgx.Query(
		ctx,
		rawQuery,
		dto.Id,
	)
	if err != nil {
		return nil, fmt.Errorf("get criterias: %w", err)
	}
	defer rows.Close()

	var criteria []AchievementCriteria

	for rows.Next() {
		var criterion AchievementCriteria

		err = rows.Scan(
			&criterion.AchievementId,
			&criterion.StatId,
			&criterion.TargetValue,
			&criterion.CriteriaText,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		criteria = append(criteria, criterion)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return criteria, nil
}

func (r *CriteriaPostgresRepository) Update(ctx context.Context, dto UpdateCriteriaDTO) (AchievementCriteria, error) {
	query := `UPDATE 
				achievement_criteria
			SET
				stat_id = $2,
				target_value = $3,
				criteria_text = $4
			WHERE
				id = $1
			RETURNING 
				achievement_id,
				stat_id,
				target_value,
				criteria_text;`

	var criteria AchievementCriteria

	err := r.pgx.QueryRow(
		ctx,
		query,
		dto.CriteriaID,
		dto.StatId,
		dto.TargetValue,
		dto.CriteriaText,
	).Scan(
		&criteria.AchievementId,
		&criteria.StatId,
		&criteria.TargetValue,
		&criteria.CriteriaText,
	)

	if err != nil {
		return AchievementCriteria{}, fmt.Errorf("scanning row: %w", err)
	}

	return criteria, nil
}

func (r *CriteriaPostgresRepository) Delete(ctx context.Context, dto DeleteCriteriaDTO) error {
	rawQuery := `DELETE FROM
					achievement_criteria
				WHERE
					id = $1;`

	_, err := r.pgx.Exec(
		ctx,
		rawQuery,
		dto.CriteriaID,
	)

	if err != nil {
		return fmt.Errorf("delete criteria: %w", err)
	}

	return nil
}
