package achievements

type CreateAchievementRequest struct {
	Key         string `json:"key" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type GetAchievementRequest struct {
	Id int `uri:"id" binding:"required"`
}

type UpdateAchievementRequest struct {
	Id          int    `uri:"id" binding:"required"`
	Key         string `json:"key" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type DeleteAchievementRequest struct {
	Id int `uri:"id" binding:"required"`
}

type AchievementResponse struct {
	Id          int    `json:"id"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
