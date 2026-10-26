package achievements

type ProgressId struct {
	Id int `uri:"id" binding:"required"`
}

type CreateProgressRequest struct {
	UserId       string `json:"user_id" binding:"required"`
	CriterionId  int    `json:"criteria_id" binding:"required"`
	CurrentValue int    `json:"current_value" binding:"required"`
}

type GetProgressRequest struct {
	ProgressId
}

type ProgressUpdateData struct {
	CurrentValue int `json:"current_value" binding:"required"`
}

type UpdateProgressRequest struct {
	ProgressId
	ProgressUpdateData
}

type DeleteProgressRequest struct {
	ProgressId
}

type GetProgressOfUserRequest struct {
	UserId        string `form:"userid" binding:"required"`
	AchievementId int    `form:"achievementid" binding:"required"`
}
