package progression

import (
	"context"
	"time"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostgresRepository struct {
	logger *zap.Logger
	pgx    *pgxpool.Pool
}

func NewPostgresRepositoryFromPgxPool(connection *pgxpool.Pool, l *zap.Logger) *PostgresRepository {
	return &PostgresRepository{
		logger: l,
		pgx:    connection,
	}
}

func (r *PostgresRepository) Update(ctx context.Context, dto UpdateDTO) (Progression, error) {
	rawQuery := `INSERT INTO
					player_progression(user_id, stat_id, current_value, updated_at)
				VALUES ($1, $2, $3, $4)
				ON CONFLICT (user_id, stat_id)
				DO UPDATE SET
					current_value = EXCLUDED.current_value,
					updated_at = EXCLUDED.updated_at
				RETURNING
					user_id,
					stat_id,
					current_value,
					updated_at;`

	now := time.Now()

	var progression Progression

	err := r.pgx.QueryRow(
		ctx,
		rawQuery,
		dto.UserId,
		dto.StatId,
		dto.CurrentValue,
		now,
	).Scan(
		&progression.UserId,
		&progression.StatId,
		&progression.CurrentValue,
		&progression.UpdatedAt,
	)

	if err != nil {
		r.logger.Error(
			"Failed to execute progression update query",
			zap.Error(err),
		)

		return Progression{}, common.TranslatePostgresError(err, r.logger)
	}

	return progression, nil
}
