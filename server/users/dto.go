package users

import "github.com/google/uuid"

type UserId struct {
	Id string `uri:"id" binding:"required,uuid"`
}

type DeleteType struct {
	Type string `form:"type" binding:"required,oneof=soft hard"`
}

type GetDTO struct {
	UserId
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

type FieldUpdate struct {
	Field    string `json:"field" binding:"required"`
	DataType string `json:"datatype" binding:"required"`
	Value    string `json:"value" binding:"required"`
}

type FieldUpdates struct {
	Updates []FieldUpdates `json:"updates" binding:"required"`
}

type UpdateDTO struct {
	UserId
	FieldUpdates
}

type ReplacementData struct {
	Username    string `json:"username" binding:"required,gte=5,lte=20"`
	DisplayName string `json:"display_name" binding:"required,gte=5,lte=10"`
	Email       string `json:"email" binding:"required,email"`
	Country     string `json:"country" binding:"required"`
	State       string `json:"state" binding:"required"`
	LangTag     string `json:"lang_tag" binding:"required"`
	Timezone    string `json:"time_zone" binding:"required"`
}

type ReplaceDTO struct {
	UserId
	ReplacementData
}

type DeleteDTO struct {
	UserId
	DeleteType
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
