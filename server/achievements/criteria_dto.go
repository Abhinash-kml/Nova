package achievements

type AchievementId struct {
	Id int `uri:"achievement_id" json:"achievement_id" binding:"required"`
}

type AchievementCriteriaId struct {
	CriteriaID int `uri:"id" binding:"required"`
}

type GetAllCriteriaDTO struct {
	Cursor int `form:"cursor"`
	Limit  int `form:"limit"`
}

type GetCriteriaDTO struct {
	AchievementCriteriaId
}

type GetCriteriaByAchievementDTO struct {
	AchievementId
}

type CreateCriteriaDTO struct {
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

type UpdateCriteriaDTO struct {
	AchievementCriteriaId
	UpdateCriteriaData
}

type DeleteCriteriaDTO struct {
	AchievementCriteriaId
}
