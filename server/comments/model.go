package comments

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	Id        uuid.UUID `json:"id" redis:"id"`
	PostId    uuid.UUID `json:"post_id" redis:"post_id"`
	AuthorId  uuid.UUID `json:"author_id" redis:"author_id"`
	Body      string    `json:"body" redis:"body"`
	CreatedAt time.Time `json:"created_at" redis:"created_at"`
	UpdatedAt time.Time `json:"updated_at" redis:"updated_at"`
}

func New(id, postid, authorid uuid.UUID, body string) Comment {
	return Comment{
		Id:       id,
		PostId:   postid,
		AuthorId: authorid,
		Body:     body,
	}
}

func (p *Comment) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

func (Comment) Unmarshall(b []byte) (Comment, error) {
	var t Comment
	err := json.Unmarshal(b, &t)
	return t, err
}
