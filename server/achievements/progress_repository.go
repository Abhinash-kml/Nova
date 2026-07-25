package achievements

import "context"

type ProgressRepository interface {
	Create(ctx context.Context, dto CreateProgressDTO) (AchievementProgress, error)
	Get(ctx context.Context, dto GetProgressDTO) (AchievementProgress, error)
	GetOfUser(ctx context.Context, dto GetProgressOfUserDTO) ([]AchievementProgress, error)
	Update(ctx context.Context, dto UpdateProgressDTO) error
	Delete(ctx context.Context, dto DeleteProgressDTO) (ProgressId, error)
}
