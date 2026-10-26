package achievements

type AchievementId struct {
	Id int `uri:"achievement_id" json:"achievement_id" binding:"required"`
}

type AchievementCriteriaId struct {
	CriteriaID int `uri:"id" binding:"required"`
}

type GetAllCriteriaRequest struct {
	Cursor int `form:"cursor,default=0" binding:"required"`
	Limit  int `form:"limit,default=10" binding:"required"`
}

type GetCriteriaRequest struct {
	AchievementCriteriaId
}

type GetCriteriaByAchievementRequest struct {
	AchievementId
}

type CreateCriteriaRequest struct {
	AchievementId int    `json:"achievement_id" binding:"required"`
	StatId        int    `json:"stat_id" binding:"required"`
	TargetValue   int    `json:"target_value" binding:"required"`
	CriteriaText  string `json:"criteria_text" binding:"required"`
}

type UpdateCriteriaData struct {
	AchievementId int    `json:"achievement_id" binding:"required"`
	StatId        int    `json:"stat_id" binding:"required"`
	TargetValue   int    `json:"target_value" binding:"required"`
	CriteriaText  string `json:"criteria_text" binding:"required"`
}

type UpdateCriteriaRequest struct {
	AchievementCriteriaId
	UpdateCriteriaData
}

type DeleteCriteriaRequest struct {
	AchievementCriteriaId
}
