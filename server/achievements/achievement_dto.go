package achievements

type CreateAchievementDTO struct {
	Key         string `json:"key" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type GetAchievementDTO struct {
	Id int `uri:"id" binding:"required"`
}

type UpdateAchievementDTO struct {
	Id          int    `uri:"id" binding:"required"`
	Key         string `json:"key" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type DeleteAchievementDTO struct {
	Id int `uri:"id" binding:"required"`
}

type AchievementResponseDTO struct {
	Id          int    `json:"id"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
