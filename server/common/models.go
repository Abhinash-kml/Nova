package common

import "github.com/google/uuid"

type BulkOpResult struct {
	Id      uuid.UUID
	Success bool
	Status  int
	Message string
}

type BulkSummary struct {
	TotalProcessed int `json:"total_processed"`
	TotalSuccess   int `json:"total_success"`
	TotalErrors    int `json:"total_errors"`
}

type OpResult struct {
	Id      uuid.UUID `json:"id"`
	Status  int       `json:"status"`
	Message string    `json:"message"`
}

type AddResult = OpResult
type DeleteResult = OpResult
type ReplaceResult = OpResult

type BulkResponse struct {
	Success  bool            `json:"success"`
	Summary  BulkSummary     `json:"summary"`
	Added    []AddResult     `json:"added,omitempty"`
	Deleted  []DeleteResult  `json:"deleted,omitempty"`
	Replaced []ReplaceResult `json:"replaced,omitempty"`
}
