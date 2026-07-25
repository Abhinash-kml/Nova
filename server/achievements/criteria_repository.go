package achievements

import "context"

type CriteriaRepository interface {
	GetAll(ctx context.Context, dto GetAllCriteriaDTO) ([]AchievementCriteria, error)
	Add(ctx context.Context, dto CreateCriteriaDTO) (AchievementCriteria, error)
	Get(ctx context.Context, dto GetCriteriaDTO) (AchievementCriteria, error)
	Update(ctx context.Context, dto UpdateCriteriaDTO) (AchievementCriteria, error)
	Delete(ctx context.Context, dto DeleteCriteriaDTO) error
	GetByAchievement(ctx context.Context, dto GetCriteriaByAchievementDTO) ([]AchievementCriteria, error)
}
