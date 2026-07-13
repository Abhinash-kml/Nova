package progression

type UserId struct {
	Id string `uri:"userid" binding:"required,uuid"`
}

type ProgressionId struct {
	StatId int `uri:"stat_id" json:"stat_id" binding:"required"`
}

type UpdateDTO struct {
	UserId
	ProgressionId
	CurrentValue int `json:"current_value" binding:"required"`
}
