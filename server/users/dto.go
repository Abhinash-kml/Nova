package users

import "github.com/google/uuid"

type UserCreateDTO struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Country     string `json:"country"`
	State       string `json:"state"`
	LangTag     string `json:"lang_tag"`
	Timezone    string `json:"time_zone"`
}

func NewUserCreateDTO(user *User) UserCreateDTO {
	return UserCreateDTO{
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		Country:     user.Country,
		State:       user.State,
		LangTag:     user.LangTag,
	}
}

type UserUpdateDTO struct {
	Id       uuid.UUID `json:"id"`
	Field    string    `json:"field"`
	DataType string    `json:"datatype"`
	Value    string    `json:"value"`
}

type UserReplaceDTO struct {
	Id          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email"`
	Country     string    `json:"country"`
	State       string    `json:"state"`
	LangTag     string    `json:"lang_tag"`
	Timezone    string    `json:"time_zone"`
}

type UserDeleteDTO struct {
	Id   uuid.UUID `json:"id"`
	Type string    `json:"type"` // Soft (disable) - Hard (delete)
}
