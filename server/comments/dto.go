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

type CommentId struct {
	Id string `uri:"id" binding:"required,uuid"`
}

type CommentBody struct {
	Body string `json:"body" binding:"required"`
}

type Body struct {
	Body string `json:"body" binding:"required"`
}

type UpdateDTO struct {
	CommentId
	Body
}

type ReplaceDTO struct {
	CommentId
	Body
}

type DeleteOptions struct {
	Type string `form:"type" binding:"required,oneof=soft hard"` // Soft (disable) - Hard (delete)
}

type DeleteDTO struct {
	CommentId
	DeleteOptions
}

type BulkCreateDTO struct {
	Comments []CreateDTO `json:"comments" binding:"required"`
}

type BulkModifyDTO struct {
	Updates []UpdateDTO `json:"updates" binding:"required"`
}

type BulkDeleteDTO struct {
	Comments []uuid.UUID `json:"comments" binding:"required"`
}
