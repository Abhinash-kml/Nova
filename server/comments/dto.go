package comments

import "github.com/google/uuid"

type GetDTO struct {
	Id string `uri:"id" binding:"required,uuid"`
}

type GetAllDTO struct {
	Cursor string `form:"cursor" binding:"required"`
	Limit  int    `form:"limit" binding:"required,gte=10,lte=20"`
}

type CreateDTO struct {
	PostId   uuid.UUID `json:"post_id" binding:"required"`
	AuthorId uuid.UUID `json:"author_id" binding:"required"`
	Body     string    `json:"body" binding:"required"`
}

type UpdateDTO struct {
	Id   string `uri:"id" binding:"required,uuid"`
	Body string `json:"body" binding:"required"`
}

type ReplaceDTO struct {
	Id   string `uri:"id" binding:"required,uuid"`
	Body string `json:"body" binding:"required"`
}

type DeleteDTO struct {
	Id   string `uri:"id" binding:"required,uuid"`
	Type string `form:"type" binding:"required,oneof=soft hard"` // Soft (disable) - Hard (delete)
}
