package posts

import "github.com/google/uuid"

type GetDTO struct {
	Id string `uri:"id" binding:"required,uuid"`
}

type GetAllDTO struct {
	Cursor        string `form:"cursor,default=AAAAAAAAAAAAAAAAAAAAAA" binding:"required"`
	Limit         int    `form:"limit,default=10" binding:"required,gte=10,lte=20,numeric"`
	SortDirection string `form:"sort,default=asc" binding:"required,oneof=asc desc"`
	SortBy        string `form:"sortby,default=id" binding:"required"`
}

type CreateDTO struct {
	Title    string    `json:"title" binding:"required,min=5,max=20"`
	Body     string    `json:"body" binding:"required,min=5,max=200"`
	AuthorId uuid.UUID `json:"author_id" binding:"required,uuid"`
}

type PostId struct {
	Id string `uri:"id" binding:"required,uuid"`
}

type FieldUpdate struct {
	Field    string `json:"field" binding:"required"`
	DataType string `json:"datatype" binding:"required"`
	Value    string `json:"value" binding:"required"`
}

type FieldUpdates struct {
	Updates []FieldUpdate `json:"updates" binding:"required"`
}

type UpdateDTO struct {
	PostId
	Title *string `json:"title"`
	Body  *string `json:"body"`
}

type ReplacementData struct {
	Title string `json:"title" binding:"required"`
	Body  string `json:"body" binding:"required"`
}

type ReplaceDTO struct {
	PostId
	ReplacementData
}

type DeleteOptions struct {
	Type string `form:"type" binding:"required,oneof=soft hard"`
}

type DeleteDTO struct {
	PostId
	DeleteOptions
}

type BulkCreateDTO struct {
	Posts []CreateDTO `json:"posts" binding:"required"`
}

type BulkModifyDTO struct {
	Updates []UpdateDTO `json:"updates" binding:"required"`
}

type BulkDeleteDTO struct {
	Posts []uuid.UUID `json:"posts" binding:"required"`
}
