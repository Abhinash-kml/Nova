package users

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id          uuid.UUID `json:"id" redis:"id"`
	Username    string    `json:"username" redis:"username"`
	DisplayName string    `json:"display_name" redis:"display_name"`
	Email       string    `json:"email" redis:"email"`
	Country     string    `json:"country" redis:"country"`
	State       string    `json:"state" redis:"state"`
	AvatarURL   string    `json:"avatar_url" redis:"avatar_url"`
	LangTag     string    `json:"lang_tag" redis:"lang_tag"`
	Timezone    string    `json:"timezone" redis:"timezone"`
	CreatedAt   time.Time `json:"created_at" redis:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" redis:"updated_at"`
	VerifiedAt  time.Time `json:"verified_at" redis:"verified_at"`
}

func New(id uuid.UUID, username, displayname, email, country, state, langtag, timezone string) User {
	return User{
		Id:          id,
		Username:    username,
		DisplayName: displayname,
		Email:       email,
		Country:     country,
		State:       state,
		LangTag:     langtag,
		Timezone:    timezone,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (u *User) Marshal() ([]byte, error) {
	return json.Marshal(u)
}

func (User) Unmarshall(b []byte) (User, error) {
	var t User
	err := json.Unmarshal(b, &t)
	return t, err
}
