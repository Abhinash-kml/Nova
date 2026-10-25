package posts

import "github.com/google/uuid"

type PostId struct {
	Id string `uri:"id" binding:"required,uuid"`
}

type GetRequest struct {
	PostId
}

type GetAllRequest struct {
	Cursor        string `form:"cursor,default=AAAAAAAAAAAAAAAAAAAAAA" binding:"required"`
	Limit         int    `form:"limit,default=10" binding:"required,gte=10,lte=20,numeric"`
	SortDirection string `form:"sort,default=asc" binding:"required,oneof=asc desc"`
	SortBy        string `form:"sortby,default=id" binding:"required"`
}

type CreateRequest struct {
	Title    string    `json:"title" binding:"required,min=5,max=20"`
	Body     string    `json:"body" binding:"required,min=5,max=200"`
	AuthorId uuid.UUID `json:"author_id" binding:"required,uuid"`
}

type UpdateRequest struct {
	PostId
	Title *string `json:"title"`
	Body  *string `json:"body"`
}

type ReplacementData struct {
	Title string `json:"title" binding:"required"`
	Body  string `json:"body" binding:"required"`
}

type ReplaceRequest struct {
	PostId
	ReplacementData
}

type DeleteOptions struct {
	Type string `form:"type" binding:"required,oneof=soft hard"`
}

type DeleteRequest struct {
	PostId
	DeleteOptions
}

type BulkCreateRequest struct {
	Posts []CreateRequest `json:"posts" binding:"required"`
}

type BulkModifyRequest struct {
	Updates []UpdateRequest `json:"updates" binding:"required"`
}

type BulkDeleteRequest struct {
	Posts []uuid.UUID `json:"posts" binding:"required"`
}
