package posts

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const PostPrefix = "post:"

type Post struct {
	Id        uuid.UUID `json:"id" redis:"id"`
	Title     string    `json:"title" redis:"title"`
	Body      string    `json:"body" redis:"body"`
	AuthorId  uuid.UUID `json:"author_id" redis:"author_id"`
	Likes     int       `json:"likes" redis:"likes"`
	Comments  int       `json:"comments" redis:"comments"`
	CreatedAt time.Time `json:"created_at" redis:"created_at"`
	UpdatedAt time.Time `json:"updated_at" redis:"updated_at"`
}

func New(id, authorid uuid.UUID, title, body string) Post {
	return Post{
		Id:        id,
		Title:     title,
		Body:      body,
		AuthorId:  authorid,
		Likes:     0,
		Comments:  0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (p *Post) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

func (Post) Unmarshall(b []byte) (Post, error) {
	var t Post
	err := json.Unmarshal(b, &t)
	return t, err
}

func (c *Post) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

// UnmarshalBinary deserializes the Redis blob back into your struct
func (c *Post) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, c)
}
