package comments

import "github.com/google/uuid"

type CommentCreateDTO struct {
	PostId   uuid.UUID `json:"post_id"`
	AuthorId uuid.UUID `json:"author_id"`
	Body     string    `json:"body"`
}

func NewCommentcreateDTO() CommentCreateDTO {
	return CommentCreateDTO{}
}

type CommentUpdateDTO struct {
	Id   uuid.UUID `json:"id"`
	Body string    `json:"body"`
}
