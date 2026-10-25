package comments

import "github.com/google/uuid"

type CommentId struct {
	Id string `uri:"id" binding:"required,uuid"`
}

type GetRequest struct {
	CommentId
}

type GetAllRequest struct {
	Cursor string `form:"cursor" binding:"required"`
	Limit  int    `form:"limit" binding:"required,gte=10,lte=20"`
}

type CreateRequest struct {
	PostId   uuid.UUID `json:"post_id" binding:"required"`
	AuthorId uuid.UUID `json:"author_id" binding:"required"`
	Body     string    `json:"body" binding:"required"`
}

type CommentBody struct {
	Body string `json:"body" binding:"required"`
}

type Body struct {
	Body string `json:"body" binding:"required"`
}

type UpdateRequest struct {
	CommentId
	Body
}

type ReplaceRequest struct {
	CommentId
	Body
}

type DeleteOptions struct {
	Type string `form:"type" binding:"required,oneof=soft hard"` // Soft (disable) - Hard (delete)
}

type DeleteRequest struct {
	CommentId
	DeleteOptions
}

type BulkCreateRequest struct {
	Comments []CreateRequest `json:"comments" binding:"required"`
}

type BulkModifyRequest struct {
	Updates []UpdateRequest `json:"updates" binding:"required"`
}

type BulkDeleteRequest struct {
	Comments []uuid.UUID `json:"comments" binding:"required"`
}
