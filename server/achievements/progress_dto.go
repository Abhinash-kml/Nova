package achievements

type ProgressId struct {
	Id int `uri:"id" binding:"required"`
}

type CreateProgressDTO struct {
	UserId       string `json:"user_id" binding:"required"`
	CriterionId  int    `json:"criteria_id" binding:"required"`
	CurrentValue int    `json:"current_value" binding:"required"`
}

type GetProgressDTO struct {
	ProgressId
}

type ProgressUpdateData struct {
	CurrentValue int `json:"current_value" binding:"required"`
}

type UpdateProgressDTO struct {
	ProgressId
	ProgressUpdateData
}

type DeleteProgressDTO struct {
	ProgressId
}

type GetProgressOfUserDTO struct {
	UserId        string `form:"userid" binding:"required"`
	AchievementId int    `form:"achievementid" binding:"required"`
}
