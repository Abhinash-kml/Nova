package achievements

import "context"

type Repository interface {
	Create(ctx context.Context, dto CreateAchievementRequest) (AchievementResponse, error)
	Get(ctx context.Context, dto GetAchievementRequest) (AchievementResponse, error)
	Update(ctx context.Context, dto UpdateAchievementRequest) (AchievementResponse, error)
	Delete(ctx context.Context, dto DeleteAchievementRequest) error

	CompleteAchievement(ctx context.Context, req CompleteAchievementRequest) error

	GetProgress(ctx context.Context, req GetProgressRequest) (AchievementProgress, error)
	UpdateProgress(ctx context.Context, req UpdateProgressRequest) error
	DeleteProgress(ctx context.Context, req DeleteProgressRequest) (ProgressId, error)

	GetProgressOfUser(ctx context.Context, req GetProgressOfUserRequest) ([]AchievementProgress, error)
	GetCompletedAchievementsOfUser(ctx context.Context, req CompletedAchievementRequest) ([]int, error)
	DeleteAllProgressOfUser(ctx context.Context, req DeleteProgressOfUserRequest) error

	GetAchievementsForEvaluator(ctx context.Context) ([]Achievement, error)
}
