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
	Name      string `json:"name"`
	Type      string `json:"type"`
	CreatedBy string `json:"created_by"`
}

type DeleteDTO struct {
	LeaderboardId
}

// To be finalized
type ModifyDTO struct {
	LeaderboardId
}

type Score struct {
	Id    uuid.UUID `json:"id"`
	Score uint      `json:"score"`
}

type GetScoreDTO struct {
	LeaderboardId
}

type ScoreDTO struct {
	Scores []Score `json:"scores"`
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
