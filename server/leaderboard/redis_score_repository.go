package leaderboard

import (
	"context"
	"fmt"

	"github.com/abhinash-kml/nova/server/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisScoreRepository struct {
	config  *config.Config
	logger  *zap.Logger
	rclient *redis.Client
}

func NewRedisScoreRepository(c *config.Config, l *zap.Logger, r *redis.Client) *RedisScoreRepository {
	return &RedisScoreRepository{
		config:  c,
		rclient: r,
		logger:  l,
	}
}

func (r *RedisScoreRepository) GetScore(ctx context.Context, dto GetScoreDTO) (ScoreDTO, error) {
	result, err := r.rclient.ZRevRangeWithScores(ctx, dto.Id, 0, int64(dto.Limit)).Result()
	if err != nil {
		return ScoreDTO{}, err
	}

	scores := ScoreDTO{
		Scores: make([]Score, len(result)),
	}

	for i := range result {
		scores.Scores[i] = Score{
			UserId: UserId{
				Id: result[i].Member.(string),
			},
			Score: uint(result[i].Score),
		}
	}

	return scores, nil
}

func (r *RedisScoreRepository) UpdateScore(ctx context.Context, dto UpdateScoreDTO) error {
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
				return fmt.Errorf("update score with best aggregation: %w", err)
			}
		} else {
			_, err = r.rclient.ZAdd(ctx, dto.Id, members...).Result()
		}

		if err != nil {
			return fmt.Errorf("update score best/set: %w", err)
		}
	}

	// For increment / decrement
	// Create transaction pipeline
	pipe := r.rclient.Pipeline()
	var err error

	for i := range dto.Scores {
		member := dto.Scores[i].Id
		score := float64(dto.Scores[i].Score)

		switch dto.AggregateType {
		case "incr":
			_, err = pipe.ZIncrBy(ctx, dto.Id, score, member).Result()
		case "decr":
			_, err = pipe.ZIncrBy(ctx, dto.Id, -score, member).Result()
		}

		if err != nil {
			return fmt.Errorf("update score incr/decr: %w", err)
		}
	}

	// Execute pipeline across the network in exactly one transaction block
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("score update pipeline: %w", err)
	}

	return nil
}

func (r *RedisScoreRepository) DeleteScore(ctx context.Context, dto DeleteScoreDTO) error {
	leaderboardId := dto.LeaderboardId.Id
	userId := dto.UserId.Id
	_, err := r.rclient.ZRem(ctx, leaderboardId, userId).Result()
	if err != nil {
		return fmt.Errorf("delete score: %w", err)
	}

	return nil
}

func (r *RedisScoreRepository) CreateNewAndLoad(ctx context.Context, leaderbaord Leaderboard, scores []Score) error {

	return nil
}
