package users

import "github.com/google/uuid"

type GetDTO struct {
	Id string `uri:"id" binding:"required,uuid"`
}

type GetAllDTO struct {
	Cursor string `form:"cursor" binding:"required"`
	Limit  int    `form:"limit" binding:"required,gte=10,lte=20"`
}

type CreateDTO struct {
	Username    string `json:"username" binding:"required,gte=5,lte=20"`
	DisplayName string `json:"display_name" binding:"required,gte=5,lte=20"`
	Email       string `json:"email" binding:"required,email"`
	Country     string `json:"country" binding:"required"`
	State       string `json:"state" binding:"required"`
	LangTag     string `json:"lang_tag" binding:"required"`
	Timezone    string `json:"time_zone" binding:"required"`
}

func NewUserCreateDTO(user *User) CreateDTO {
	return CreateDTO{
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		Country:     user.Country,
		State:       user.State,
		LangTag:     user.LangTag,
	}
}

type FieldUpdates struct {
	Field    string `json:"field" binding:"required"`
	DataType string `json:"datatype" binding:"required"`
	Value    string `json:"value" binding:"required"`
}

type UpdateDTO struct {
	Id      string         `uri:"id" binding:"required,uuid"`
	Updates []FieldUpdates `json:"updates" binding:"required"`
}

type ReplaceDTO struct {
	Id          string `uri:"id" binding:"required,uuid"`
	Username    string `json:"username" binding:"required,gte=5,lte=20"`
	DisplayName string `json:"display_name" binding:"required,gte=5,lte=10"`
	Email       string `json:"email" binding:"required,email"`
	Country     string `json:"country" binding:"required"`
	State       string `json:"state" binding:"required"`
	LangTag     string `json:"lang_tag" binding:"required"`
	Timezone    string `json:"time_zone" binding:"required"`
}

type DeleteDTO struct {
	Id   string `uri:"id" binding:"required,uuid"`
	Type string `form:"type" binding:"required,oneof=soft hard"` // Soft (disable) - Hard (delete)
}

type UserId struct {
	Id string `uri:"id" binding:"required,uuid"`
}

type DeleteType struct {
	Type string `form:"type" binding:"required,oneof=soft hard"`
}

type BulkCreateDTO struct {
	Users []CreateDTO `json:"users" binding:"required"`
}

type BulkModifyDTO struct {
	Updates []UpdateDTO `json:"updates" binding:"required"`
}

type BulkDeleteDTO struct {
	Users []uuid.UUID `json:"users" binding:"required"`
}
