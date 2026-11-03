package stats

type StatsId struct {
	Id int `uri:"id" json:"id" binding:"required"`
}

type GetRequest struct {
	StatsId
}

type GetAllRequest struct {
	Cursor int `form:"cursor,default=1" binding:"required"`
	Limit  int `form:"limit,default=10" binding:"required,gte=10,lte=20"`
}

type CreateRequest struct {
	Key         string `json:"key" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	StartValue  int    `json:"start_value" binding:"required"`
}

type UpdateRequest struct {
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

type ReplaceRequest struct {
	StatsId
	ReplacementData
}

type DeleteRequest struct {
	StatsId
}

type UserId struct {
	Id string `uri:"id" json:"id" binding:"required"`
}

type GetPlayerStatRequest struct {
	UserId
}

type StatResponse struct {
	Id    int64 `json:"id"`
	Value int   `json:"value"`
}

type PlayerStatsResponse struct {
	Stats []StatResponse `json:"player_stats"`
}

type DeletePlayerStatRequest struct {
	UserId
}

type DeletePlayerStatSpecificRequest struct {
	UserId
	StatsId string `form:"statid" binding:"required,uuid"`
}

type IncomingPlayerStat struct {
	StatId string `json:"stat_id"`
	Value  string `json:"value"`
}

type UpdatePlayerStatRequest struct {
	UserId
	IncomingPlayerStat
}
