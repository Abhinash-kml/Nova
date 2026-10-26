package achievements

import "context"

type Repository interface {
	Create(ctx context.Context, dto CreateAchievementRequest) (AchievementResponse, error)
	Get(ctx context.Context, dto GetAchievementRequest) (AchievementResponse, error)
	Update(ctx context.Context, dto UpdateAchievementRequest) error
	Delete(ctx context.Context, dto DeleteAchievementRequest) error
}
