package stats

type StatsId struct {
	Id int `uri:"id" json:"id" binding:"required"`
}

type GetDTO struct {
	StatsId
}

type GetAllDTO struct {
	Cursor string `form:"cursor" binding:"required"`
	Limit  int    `form:"limit" binding:"required,gte=10,lte=20" default:"10"`
}

type CreateDTO struct {
	Key         string `json:"key" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	StartValue  int    `json:"start_value" binding:"required"`
}

type UpdateDTO struct {
	StatsId

	Name       *string `json:"name"`
	StartValue *int    `json:"start_value"`
}

type ReplacementData struct {
	Name        string `json:"name" binding:"required"`
	Key         string `json:"key" binding:"required"`
	Description string `json:"description" binding:"required"`
	StartValue  int    `json:"start_value" binding:"required"`
}

type ReplaceDTO struct {
	StatsId
	ReplacementData
}

type DeleteDTO struct {
	StatsId
}

type UserId struct {
	Id string `uri:"id" json:"id" binding:"required"`
}

type GetPlayerStatDTO struct {
	UserId
}

type StatResponseDTO struct {
	Id    int64 `json:"id"`
	Value int   `json:"value"`
}

type PlayerStatsResponseDTO struct {
	Stats []StatResponseDTO `json:"player_stats"`
}

type DeletePlayerStatDTO struct {
	UserId
}

type DeletePlayerStatSpecificDTO struct {
	UserId
	StatsId string `form:"statid" binding:"required,uuid"`
}

type IncomingPlayerStat struct {
	StatId string `json:"stat_id"`
	Value  int    `json:"value"`
}

type UpdatePlayerStatDTO struct {
	UserId
	IncomingPlayerStat
}
