package stats

type StatsId struct {
	Id int `uri:"id" json:"id" binding:"required"`
}

type GetDTO struct {
	StatsId
}

type GetAllDTO struct {
	Cursor int `form:"cursor" binding:"required,gte=0"`
	Limit  int `form:"limit" binding:"required,gte=10,lte=20"`
}

type CreateDTO struct {
	Name       string `json:"name" binding:"required"`
	StartValue int    `json:"startvalue" binding:"required"`
}

type UpdateDTO struct {
	StatsId

	Name       *string `json:"name"`
	StartValue *int    `json:"startvalue"`
}

type ReplacementData struct {
	Name       string `json:"name" binding:"required"`
	StartValue int    `json:"startvalue" binding:"required"`
}

type ReplaceDTO struct {
	StatsId
	ReplacementData
}

type DeleteDTO struct {
	StatsId
}
