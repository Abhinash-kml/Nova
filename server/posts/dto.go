package posts

import "github.com/google/uuid"

type GetDTO struct {
	Id string `uri:"id" binding:"required,uuid"`
}

type GetAllDTO struct {
	Cursor string `form:"cursor" binding:"required"`
	Limit  int    `form:"limit" binding:"required,gte=10,lte=20"`
}

type CreateDTO struct {
	Title    string    `json:"title" binding:"required,min=5,max=20"`
	Body     string    `json:"body" binding:"required,min=5,max=200"`
	AuthorId uuid.UUID `json:"author_id" binding:"required,uuid"`
}

type UpdateDTO struct {
	Id       string `uri:"id" binding:"required"`
	Field    string `json:"field" binding:"required"`
	DataType string `json:"datatype" binding:"required"`
	Value    string `json:"value" binding:"required"`
}

type ReplaceDTO struct {
	Id    uuid.UUID `uri:"id" binding:"required"`
	Title string    `json:"title" binding:"required"`
	Body  string    `json:"body" binding:"required"`
}

type DeleteDTO struct {
	Id   uuid.UUID `uri:"id" binding:"required"`
	Type string    `form:"type" binding:"required,oneof=soft hard"` // Soft (disable) - Hard (delete)
}
