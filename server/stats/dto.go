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
	StartValue  int    `json:"startvalue" binding:"required"`
}

type UpdateDTO struct {
	StatsId

	Name       *string `json:"name"`
	StartValue *int    `json:"startvalue"`
}

type ReplacementData struct {
	Name        string `json:"name" binding:"required"`
	Key         string `json:"key" binding:"required"`
	Description string `json:"description" binding:"required"`
	StartValue  int    `json:"startvalue" binding:"required"`
}

type ReplaceDTO struct {
	StatsId
	ReplacementData
}

type DeleteDTO struct {
	StatsId
}
