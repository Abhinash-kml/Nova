package achievements

import "context"

type CompletedAchievementRepository interface {
	Create(ctx context.Context, req CompleteAchievementRequest) error
}
