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
	Leaderboard
}

type Score struct {
	Id    uuid.UUID `json:"id"`
	Score uint      `json:"score"`
}

type GetScoreDTO struct {
	LeaderboardId
}

type LeaderboardScoreDTO struct {
	Scores []Score `json:"scores"`
}

type UpdateOptions struct {
	AggregateType string `uri:"type" binding:"required,oneof=incr decr best set"`
}

type UpdateScoreDTO struct {
	LeaderboardId
}

type ModifyScoreDTO struct {
}

type DeleteScoreDTO struct {
	LeaderboardId
	UserId
}
