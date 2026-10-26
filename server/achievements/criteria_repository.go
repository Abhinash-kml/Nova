package achievements

import "context"

type CriteriaRepository interface {
	GetAll(ctx context.Context, dto GetAllCriteriaRequest) ([]AchievementCriteria, error)
	Add(ctx context.Context, dto CreateCriteriaRequest) (AchievementCriteria, error)
	Get(ctx context.Context, dto GetCriteriaRequest) (AchievementCriteria, error)
	Update(ctx context.Context, dto UpdateCriteriaRequest) (AchievementCriteria, error)
	Delete(ctx context.Context, dto DeleteCriteriaRequest) error
	GetByAchievement(ctx context.Context, dto GetCriteriaByAchievementRequest) ([]AchievementCriteria, error)
}
