package leaderboard

import "github.com/google/uuid"

type LeaderboardId struct {
	Id string `uri:"id" binding:"required,uuid"`
}

type UserId struct {
	Id string `form:"id" binding:"required,uuid"`
}

type GetDTO struct {
	LeaderboardId
}

type GetAllDTO struct {
	Cursor string `form:"cursor" binding:"required"`
	Limit  int    `form:"limit" binding:"required,gte=10,lte=20"`
}

type CreateDTO struct {
	Name            string `json:"name" binding:"required"`
	Type            string `json:"type" binding:"required"`
	ProcessInterval int    `json:"process_interval" binding:"required"`
	CreatedBy       string `json:"created_by" binding:"required"`
}

type DeleteDTO struct {
	LeaderboardId
}

type LeaderboardModifications struct {
	Type            string `json:"type" binding:"required"`
	ProcessInterval int    `json:"process_interval" binding:"required"`
}

type ModifyDTO struct {
	LeaderboardId
	LeaderboardModifications
}

type Score struct {
	Id    uuid.UUID `json:"id" binding:"required"`
	Score uint      `json:"score" binding:"required"`
}

type GetScoreDTO struct {
	LeaderboardId
}

type ScoreDTO struct {
	Scores []Score `json:"scores" binding:"required"`
}

type UpdateOptions struct {
	AggregateType string `uri:"operator" binding:"required,oneof=incr decr best set"`
}

type UpdateScoreDTO struct {
	LeaderboardId
	UpdateOptions
	ScoreDTO
}

type ModifyScoreDTO struct {
}

type DeleteScoreDTO struct {
	LeaderboardId
	UserId
}
