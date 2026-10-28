package achievements

type AchievementId struct {
	Id int `uri:"id" binding:"required"`
}

type GetAchievementRequest struct {
	AchievementId
}

type CreateAchievementRequest struct {
	Key         string `json:"key" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Stat        string `json:"stat" binding:"required"`
	TargetValue int    `json:"target_value" binding:"required"`
}

type UpdateAchievementRequest struct {
	Id          int     `json:"id" binding:"required"`
	Key         *string `json:"key" binding:"required"`
	Name        *string `json:"name" binding:"required"`
	Description *string `json:"description" binding:"required"`
	Stat        *string `json:"stat" binding:"required"`
	TargetValue *int    `json:"target_value" binding:"required"`
}

type DeleteAchievementRequest struct {
	AchievementId
}

type AchievementResponse struct {
	Id          int    `json:"id"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Stat        string `json:"stat"`
	TargetValue int    `json:"targetValue"`
}

type CompleteAchievementRequest struct {
	UserId        string `json:"user_id" binding:"required"`
	AchievementId int    `json:"achievement_id" binding:"required"`
}

type ProgressId struct {
	Id int `uri:"id" binding:"required"`
}

type GetProgressRequest struct {
	ProgressId
}

type ProgressUpdateData struct {
	UserId        string `json:"user_id" binding:"required"`
	AchievementId string `json:"achievement_id" binding:"required"`
	CurrentValue  int    `json:"current_value" binding:"required"`
}

type UpdateProgressRequest struct {
	ProgressUpdateData
}

type DeleteProgressRequest struct {
	ProgressId
}

type GetProgressOfUserRequest struct {
	UserId string `form:"user_id" binding:"required"`
}

type DeleteProgressOfUserRequest struct {
	UserId string `uri:"id" binding:"required"`
}
