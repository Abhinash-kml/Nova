package achievements

import "context"

type ProgressRepository interface {
	Create(ctx context.Context, dto CreateProgressRequest) (AchievementProgress, error)
	Get(ctx context.Context, dto GetProgressRequest) (AchievementProgress, error)
	GetOfUser(ctx context.Context, dto GetProgressOfUserRequest) ([]AchievementProgress, error)
	Update(ctx context.Context, dto UpdateProgressRequest) error
	Delete(ctx context.Context, dto DeleteProgressRequest) (ProgressId, error)
}
