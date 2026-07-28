package achievements

import (
	"context"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/abhinash-kml/nova/server/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostgresCompletedAchievementRepository struct {
	pgx    *pgxpool.Pool
	config *config.Config
	logger *zap.Logger
}

func NewPostgresCompletedAchievementRepository(pool *pgxpool.Pool, c *config.Config, l *zap.Logger) *PostgresCompletedAchievementRepository {
	return &PostgresCompletedAchievementRepository{
		pgx:    pool,
		config: c,
		logger: l,
	}
}

func (r *PostgresCompletedAchievementRepository) Create(ctx context.Context, dto CompleteAchievementDTO) error {
	rawQuery := `INSERT INTO 
				completed_achievements(
					user_id,
					achievement_id
				) 
				VALUES ($1, $2);`

	_, err := r.pgx.Exec(ctx, rawQuery, dto.UserId, dto.AchievementId)
	if err != nil {
		r.logger.Error("Failed to exec add user completed achievement query", zap.Error(err))
		return common.TranslatePostgresError(err, r.logger)
	}

	return nil
}
