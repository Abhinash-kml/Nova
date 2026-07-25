package achievements

import "context"

type Repository interface {
	Create(ctx context.Context, dto CreateAchievementDTO) (AchievementResponseDTO, error)
	Get(ctx context.Context, dto GetAchievementDTO) (AchievementResponseDTO, error)
	Update(ctx context.Context, dto UpdateAchievementDTO) error
	Delete(ctx context.Context, dto DeleteAchievementDTO) error
}
