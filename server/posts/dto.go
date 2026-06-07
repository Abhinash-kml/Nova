package posts

import "github.com/google/uuid"

type PostCreateDTO struct {
	Title    string    `json:"title"`
	Body     string    `json:"body"`
	AuthorId uuid.UUID `json:"author_id"`
}

func NewPostCreateDTO() PostCreateDTO {
	return PostCreateDTO{}
}

type PostUpdateDTO struct {
	Id       uuid.UUID `json:"id"`
	Field    string    `json:"field"`
	DataType string    `json:"datatype"`
	Value    string    `json:"value"`
}

type PostReplaceDTO struct {
	Id    uuid.UUID `json:"id"`
	Title string    `json:"title"`
	Body  string    `json:"body"`
}

type PostDeleteDTO struct {
	Id   uuid.UUID `json:"id"`
	Type string    `json:"type"` // Soft (disable) - Hard (delete)
}
