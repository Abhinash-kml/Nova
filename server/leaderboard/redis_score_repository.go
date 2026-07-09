package leaderboard

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisRepository struct {
	logger  *zap.Logger
	rclient *redis.Client
}

func NewRedisReposiory(l *zap.Logger, c *redis.Client) *RedisRepository {
	return &RedisRepository{
		rclient: c,
		logger:  l,
	}
}

func (r *RedisRepository) GetScore(ctx context.Context, dto GetScoreDTO) (ScoreDTO, error) {
	result, err := r.rclient.ZRevRangeWithScores(ctx, dto.Id, 0, 99).Result()
	if err != nil {
		return ScoreDTO{}, err
	}

	scores := ScoreDTO{
		Scores: make([]Score, 0, len(result)),
	}

	for i := range result {
		id, _ := uuid.Parse(fmt.Sprintf("%v", result[i].Member))

		scores.Scores[i] = Score{
			Id:    id,
			Score: uint(result[i].Score),
		}
	}

	return scores, nil
}

func (r *RedisRepository) UpdateScore(ctx context.Context, dto UpdateScoreDTO) error {
	if dto.AggregateType == "best" || dto.AggregateType == "set" {
		members := make([]redis.Z, len(dto.Scores))
		for i := range members {
			members[i] = redis.Z{
				Member: dto.Scores[i].Id,
				Score:  float64(dto.Scores[i].Score),
			}
		}

		var err error
		if dto.AggregateType == "best" {
			_, err = r.rclient.ZAddArgs(ctx, dto.Id, redis.ZAddArgs{
				GT:      true,
				Members: members,
			}).Result()
			if err != nil {
				r.logger.Error("Failed to update score with best aggregation")
				return err
			}
		} else {
			_, err = r.rclient.ZAdd(ctx, dto.Id, members...).Result()
		}

		if err != nil {
			r.logger.Error("Failed to update score",
				zap.String("leaderboard", dto.Id),
				zap.String("aggregation", dto.AggregateType),
				zap.Int("updates", len(dto.Scores)))
		}
	}

	// For increment / decrement
	// Create transaction pipeline
	pipe := r.rclient.Pipeline()

	for i := range dto.Scores {
		member := dto.Scores[i].Id.String()
		score := float64(dto.Scores[i].Score)

		switch dto.AggregateType {
		case "incr":
			pipe.ZIncrBy(ctx, dto.Id, score, member)
		case "decr":
			pipe.ZIncrBy(ctx, dto.Id, -score, member)
		}
	}

	// Execute pipeline across the network in exactly one transaction block
	_, err := pipe.Exec(ctx)
	if err != nil {
		r.logger.Error("Pipelined incremental update failed",
			zap.String("leaderboard", dto.Id),
			zap.String("operator", dto.AggregateType),
			zap.Error(err),
		)
		return err
	}

	return nil
}

func (r *RedisRepository) DeleteScore(ctx context.Context, dto DeleteScoreDTO) error {
	leaderboardId := dto.LeaderboardId.Id
	userId := dto.UserId.Id
	_, err := r.rclient.ZRem(ctx, leaderboardId, userId).Result()
	if err != nil {
		r.logger.Error("Failed to delete score",
			zap.String("leaderboard", leaderboardId),
			zap.String("userid", userId))

		return err
	}

	return nil
}
