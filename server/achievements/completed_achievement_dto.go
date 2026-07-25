package achievements

type CompleteAchievementDTO struct {
	UserId        string `json:"user_id" binding:"required"`
	AchievementId int    `json:"achievement_id" binding:"required"`
}
