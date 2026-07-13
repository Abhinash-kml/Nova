package progression

type ProgressionId struct {
	UserId int `uri:"user_id" json:"user_id" binding:"required"`
	StatId int `uri:"stat_id" json:"stat_id" binding:"required"`
}

type UpdateDTO struct {
	ProgressionId
	CurrentValue int `json:"current_value" binding:"required"`
}
