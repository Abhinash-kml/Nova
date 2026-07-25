package achievements

import "context"

type CompletedAchievementRepository interface {
	Create(ctx context.Context, dto CompleteAchievementDTO) error
}
