package users

import (
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
	VerifiedAt  time.Time `json:"verified_at" redis:"verfied_at"`
}

func NewUser(id uuid.UUID, username, displayname, email, country, state, langtag, timezone string) *User {
	return &User{
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
